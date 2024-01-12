package common

import (
	"context"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/authzed/spicedb/internal/datastore/revisions"
	"github.com/authzed/spicedb/pkg/datastore"
	core "github.com/authzed/spicedb/pkg/proto/core/v1"
	"github.com/authzed/spicedb/pkg/tuple"
)

const (
	tuple1 = "docs:1#reader@user:1"
	tuple2 = "docs:2#editor@user:2"
)

var (
	rev1             = revisions.NewForTransactionID(1)
	rev2             = revisions.NewForTransactionID(2)
	rev3             = revisions.NewForTransactionID(3)
	revOneMillion    = revisions.NewForTransactionID(1_000_000)
	revOneMillionOne = revisions.NewForTransactionID(1_000_001)
)

func TestChanges(t *testing.T) {
	type changeEntry struct {
		revision           uint64
		relationship       string
		op                 core.RelationTupleUpdate_Operation
		deletedNamespaces  []string
		deletedCaveats     []string
		changedDefinitions []datastore.SchemaDefinition
	}

	testCases := []struct {
		name     string
		script   []changeEntry
		expected []datastore.RevisionChanges
	}{
		{
			"empty",
			[]changeEntry{},
			[]datastore.RevisionChanges{},
		},
		{
			"deleted namespace",
			[]changeEntry{
				{1, "", 0, []string{"somenamespace"}, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: nil, DeletedNamespaces: []string{"somenamespace"}},
			},
		},
		{
			"deleted caveat",
			[]changeEntry{
				{1, "", 0, nil, []string{"somecaveat"}, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: nil, DeletedCaveats: []string{"somecaveat"}},
			},
		},
		{
			"changed namespace",
			[]changeEntry{
				{1, "", 0, nil, nil, []datastore.SchemaDefinition{&core.NamespaceDefinition{
					Name: "somenamespace",
				}}},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: nil, ChangedDefinitions: []datastore.SchemaDefinition{&core.NamespaceDefinition{
					Name: "somenamespace",
				}}},
			},
		},
		{
			"create",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
			},
		},
		{
			"delete",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_DELETE, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{del(tuple1)}},
			},
		},
		{
			"in-order touch",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_DELETE, nil, nil, nil},
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
			},
		},
		{
			"reverse-order touch",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1, tuple1, core.RelationTupleUpdate_DELETE, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
			},
		},
		{
			"create and delete",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1, tuple2, core.RelationTupleUpdate_DELETE, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), del(tuple2)}},
			},
		},
		{
			"multiple creates",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1, tuple2, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), touch(tuple2)}},
			},
		},
		{
			"duplicates",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
			},
		},
		{
			"create then touch",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{2, tuple1, core.RelationTupleUpdate_DELETE, nil, nil, nil},
				{2, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
				{Revision: rev2, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
			},
		},
		{
			"big revision gap",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1_000_000, tuple1, core.RelationTupleUpdate_DELETE, nil, nil, nil},
				{1_000_000, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
				{Revision: revOneMillion, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
			},
		},
		{
			"out of order",
			[]changeEntry{
				{1_000_000, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1_000_000, tuple1, core.RelationTupleUpdate_DELETE, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
				{Revision: revOneMillion, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
			},
		},
		{
			"changed then deleted namespace",
			[]changeEntry{
				{1, "", 0, nil, nil, []datastore.SchemaDefinition{&core.NamespaceDefinition{
					Name: "somenamespace",
				}}},
				{1, "", 0, []string{"somenamespace"}, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, DeletedNamespaces: []string{"somenamespace"}},
			},
		},
		{
			"changed then deleted caveat",
			[]changeEntry{
				{1, "", 0, nil, nil, []datastore.SchemaDefinition{&core.CaveatDefinition{
					Name: "somecaveat",
				}}},
				{1, "", 0, nil, []string{"somecaveat"}, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, DeletedCaveats: []string{"somecaveat"}},
			},
		},
		{
			"deleted then changed namespace",
			[]changeEntry{
				{1, "", 0, []string{"somenamespace"}, nil, nil},
				{1, "", 0, nil, nil, []datastore.SchemaDefinition{&core.NamespaceDefinition{
					Name: "somenamespace",
				}}},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, ChangedDefinitions: []datastore.SchemaDefinition{&core.NamespaceDefinition{
					Name: "somenamespace",
				}}},
			},
		},
		{
			"deleted then changed caveat",
			[]changeEntry{
				{1, "", 0, nil, []string{"somecaveat"}, nil},
				{1, "", 0, nil, nil, []datastore.SchemaDefinition{&core.CaveatDefinition{
					Name: "somecaveat",
				}}},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, ChangedDefinitions: []datastore.SchemaDefinition{&core.CaveatDefinition{
					Name: "somecaveat",
				}}},
			},
		},
		{
			"changed namespace then deleted caveat",
			[]changeEntry{
				{1, "", 0, nil, nil, []datastore.SchemaDefinition{&core.NamespaceDefinition{
					Name: "somenamespaceorcaveat",
				}}},
				{1, "", 0, nil, []string{"somenamespaceorcaveat"}, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, DeletedCaveats: []string{"somenamespaceorcaveat"}, ChangedDefinitions: []datastore.SchemaDefinition{&core.NamespaceDefinition{
					Name: "somenamespaceorcaveat",
				}}},
			},
		},
		{
			"kitchen sink relationships",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{2, tuple1, core.RelationTupleUpdate_DELETE, nil, nil, nil},
				{1_000_000, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},

				{1, tuple2, core.RelationTupleUpdate_DELETE, nil, nil, nil},
				{2, tuple2, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1_000_000, tuple2, core.RelationTupleUpdate_DELETE, nil, nil, nil},
				{1_000_000, tuple2, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), del(tuple2)}},
				{Revision: rev2, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple2), del(tuple1)}},
				{Revision: revOneMillion, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), touch(tuple2)}},
			},
		},
		{
			"kitchen sink",
			[]changeEntry{
				{1, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{2, tuple1, core.RelationTupleUpdate_DELETE, nil, nil, nil},
				{1_000_000, tuple1, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1_000_001, "", 0, []string{"deletednamespace"}, nil, nil},

				{3, "", 0, nil, nil, []datastore.SchemaDefinition{
					&core.NamespaceDefinition{Name: "midns"},
				}},

				{1, tuple2, core.RelationTupleUpdate_DELETE, nil, nil, nil},
				{2, tuple2, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1_000_000, tuple2, core.RelationTupleUpdate_DELETE, nil, nil, nil},
				{1_000_000, tuple2, core.RelationTupleUpdate_TOUCH, nil, nil, nil},
				{1_000_001, "", 0, nil, []string{"deletedcaveat"}, nil},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), del(tuple2)}},
				{Revision: rev2, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple2), del(tuple1)}},
				{Revision: rev3, ChangedDefinitions: []datastore.SchemaDefinition{
					&core.NamespaceDefinition{Name: "midns"},
				}},
				{Revision: revOneMillion, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), touch(tuple2)}},
				{Revision: revOneMillionOne, DeletedNamespaces: []string{"deletednamespace"}, DeletedCaveats: []string{"deletedcaveat"}},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			ctx := context.Background()
			ch := NewChanges(revisions.TransactionIDKeyFunc, datastore.WatchRelationships|datastore.WatchSchema)
			for _, step := range tc.script {
				if step.relationship != "" {
					rel := tuple.MustParse(step.relationship)
					err := ch.AddRelationshipChange(ctx, revisions.NewForTransactionID(step.revision), rel, step.op)
					require.NoError(err)
				}

				for _, changed := range step.changedDefinitions {
					ch.AddChangedDefinition(ctx, revisions.NewForTransactionID(step.revision), changed)
				}

				for _, ns := range step.deletedNamespaces {
					ch.AddDeletedNamespace(ctx, revisions.NewForTransactionID(step.revision), ns)
				}

				for _, c := range step.deletedCaveats {
					ch.AddDeletedCaveat(ctx, revisions.NewForTransactionID(step.revision), c)
				}
			}

			require.Equal(
				canonicalize(tc.expected),
				canonicalize(ch.AsRevisionChanges(revisions.TransactionIDKeyLessThanFunc)),
			)
		})
	}
}

func TestFilteredSchemaChanges(t *testing.T) {
	ctx := context.Background()
	ch := NewChanges(revisions.TransactionIDKeyFunc, datastore.WatchSchema)
	require.True(t, ch.IsEmpty())

	require.NoError(t, ch.AddRelationshipChange(ctx, rev1, tuple.MustParse("document:firstdoc#viewer@user:tom"), core.RelationTupleUpdate_TOUCH))
	require.True(t, ch.IsEmpty())
}

func TestFilteredRelationshipChanges(t *testing.T) {
	ctx := context.Background()
	ch := NewChanges(revisions.TransactionIDKeyFunc, datastore.WatchRelationships)
	require.True(t, ch.IsEmpty())

	ch.AddDeletedNamespace(ctx, rev3, "deletedns3")
	require.True(t, ch.IsEmpty())
}

func TestFilterAndRemoveRevisionChanges(t *testing.T) {
	ctx := context.Background()
	ch := NewChanges(revisions.TransactionIDKeyFunc, datastore.WatchRelationships|datastore.WatchSchema)

	require.True(t, ch.IsEmpty())

	ch.AddDeletedNamespace(ctx, rev1, "deletedns1")
	ch.AddDeletedNamespace(ctx, rev2, "deletedns2")
	ch.AddDeletedNamespace(ctx, rev3, "deletedns3")

	require.False(t, ch.IsEmpty())

	results := ch.FilterAndRemoveRevisionChanges(revisions.TransactionIDKeyLessThanFunc, rev3)
	require.Equal(t, 2, len(results))
	require.False(t, ch.IsEmpty())

	require.Equal(t, []datastore.RevisionChanges{
		{
			Revision:           rev1,
			DeletedNamespaces:  []string{"deletedns1"},
			DeletedCaveats:     []string{},
			ChangedDefinitions: []datastore.SchemaDefinition{},
		},
		{
			Revision:           rev2,
			DeletedNamespaces:  []string{"deletedns2"},
			DeletedCaveats:     []string{},
			ChangedDefinitions: []datastore.SchemaDefinition{},
		},
	}, results)

	remaining := ch.AsRevisionChanges(revisions.TransactionIDKeyLessThanFunc)
	require.Equal(t, 1, len(remaining))

	require.Equal(t, []datastore.RevisionChanges{
		{
			Revision:           rev3,
			DeletedNamespaces:  []string{"deletedns3"},
			DeletedCaveats:     []string{},
			ChangedDefinitions: []datastore.SchemaDefinition{},
		},
	}, remaining)

	results = ch.FilterAndRemoveRevisionChanges(revisions.TransactionIDKeyLessThanFunc, revOneMillion)
	require.Equal(t, 1, len(results))
	require.True(t, ch.IsEmpty())

	results = ch.FilterAndRemoveRevisionChanges(revisions.TransactionIDKeyLessThanFunc, revOneMillionOne)
	require.Equal(t, 0, len(results))
	require.True(t, ch.IsEmpty())
}

func TestHLCOrdering(t *testing.T) {
	ctx := context.Background()

	ch := NewChanges(revisions.HLCKeyFunc, datastore.WatchRelationships|datastore.WatchSchema)
	require.True(t, ch.IsEmpty())

	rev1, err := revisions.HLCRevisionFromString("1.0000000001")
	require.NoError(t, err)

	rev0, err := revisions.HLCRevisionFromString("1")
	require.NoError(t, err)

	err = ch.AddRelationshipChange(ctx, rev1, tuple.MustParse("document:foo#viewer@user:tom"), core.RelationTupleUpdate_DELETE)
	require.NoError(t, err)

	err = ch.AddRelationshipChange(ctx, rev0, tuple.MustParse("document:foo#viewer@user:tom"), core.RelationTupleUpdate_TOUCH)
	require.NoError(t, err)

	remaining := ch.AsRevisionChanges(revisions.HLCKeyLessThanFunc)
	require.Equal(t, 2, len(remaining))

	require.Equal(t, []datastore.RevisionChanges{
		{
			Revision: rev0,
			RelationshipChanges: []*core.RelationTupleUpdate{
				tuple.Touch(tuple.MustParse("document:foo#viewer@user:tom")),
			},
			DeletedNamespaces:  []string{},
			DeletedCaveats:     []string{},
			ChangedDefinitions: []datastore.SchemaDefinition{},
		},
		{
			Revision: rev1,
			RelationshipChanges: []*core.RelationTupleUpdate{
				tuple.Delete(tuple.MustParse("document:foo#viewer@user:tom")),
			},
			DeletedNamespaces:  []string{},
			DeletedCaveats:     []string{},
			ChangedDefinitions: []datastore.SchemaDefinition{},
		},
	}, remaining)
}

func TestHLCSameRevision(t *testing.T) {
	ctx := context.Background()

	ch := NewChanges(revisions.HLCKeyFunc, datastore.WatchRelationships|datastore.WatchSchema)
	require.True(t, ch.IsEmpty())

	rev0, err := revisions.HLCRevisionFromString("1")
	require.NoError(t, err)

	rev0again, err := revisions.HLCRevisionFromString("1")
	require.NoError(t, err)

	err = ch.AddRelationshipChange(ctx, rev0, tuple.MustParse("document:foo#viewer@user:tom"), core.RelationTupleUpdate_TOUCH)
	require.NoError(t, err)

	err = ch.AddRelationshipChange(ctx, rev0again, tuple.MustParse("document:foo#viewer@user:sarah"), core.RelationTupleUpdate_TOUCH)
	require.NoError(t, err)

	remaining := ch.AsRevisionChanges(revisions.HLCKeyLessThanFunc)
	require.Equal(t, 1, len(remaining))

	expected := []*core.RelationTupleUpdate{
		tuple.Touch(tuple.MustParse("document:foo#viewer@user:tom")),
		tuple.Touch(tuple.MustParse("document:foo#viewer@user:sarah")),
	}
	slices.SortFunc(expected, func(i, j *core.RelationTupleUpdate) int {
		iStr := tuple.StringWithoutCaveat(i.Tuple)
		jStr := tuple.StringWithoutCaveat(j.Tuple)
		return strings.Compare(iStr, jStr)
	})

	slices.SortFunc(remaining[0].RelationshipChanges, func(i, j *core.RelationTupleUpdate) int {
		iStr := tuple.StringWithoutCaveat(i.Tuple)
		jStr := tuple.StringWithoutCaveat(j.Tuple)
		return strings.Compare(iStr, jStr)
	})

	require.Equal(t, []datastore.RevisionChanges{
		{
			Revision:            rev0,
			RelationshipChanges: expected,
			DeletedNamespaces:   []string{},
			DeletedCaveats:      []string{},
			ChangedDefinitions:  []datastore.SchemaDefinition{},
		},
	}, remaining)
}

func TestCanonicalize(t *testing.T) {
	testCases := []struct {
		name            string
		input, expected []datastore.RevisionChanges
	}{
		{
			"empty",
			[]datastore.RevisionChanges{},
			[]datastore.RevisionChanges{},
		},
		{
			"single entries",
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1)}},
			},
		},
		{
			"tuples out of order",
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{del(tuple2), touch(tuple1)}},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), del(tuple2)}},
			},
		},
		{
			"operations out of order",
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{del(tuple1), touch(tuple1)}},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), del(tuple1)}},
			},
		},
		{
			"equal entries",
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), touch(tuple1)}},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), touch(tuple1)}},
			},
		},
		{
			"already canonical",
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), del(tuple2)}},
				{Revision: rev2, RelationshipChanges: []*core.RelationTupleUpdate{del(tuple1), touch(tuple2)}},
				{Revision: revOneMillion, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), touch(tuple2)}},
			},
			[]datastore.RevisionChanges{
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), del(tuple2)}},
				{Revision: rev2, RelationshipChanges: []*core.RelationTupleUpdate{del(tuple1), touch(tuple2)}},
				{Revision: revOneMillion, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), touch(tuple2)}},
			},
		},
		{
			"revisions allowed out of order",
			[]datastore.RevisionChanges{
				{Revision: revOneMillion, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), touch(tuple2)}},
				{Revision: rev2, RelationshipChanges: []*core.RelationTupleUpdate{del(tuple1), touch(tuple2)}},
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), del(tuple2)}},
			},
			[]datastore.RevisionChanges{
				{Revision: revOneMillion, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), touch(tuple2)}},
				{Revision: rev2, RelationshipChanges: []*core.RelationTupleUpdate{del(tuple1), touch(tuple2)}},
				{Revision: rev1, RelationshipChanges: []*core.RelationTupleUpdate{touch(tuple1), del(tuple2)}},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, canonicalize(tc.input))
		})
	}
}

func touch(relationship string) *core.RelationTupleUpdate {
	return &core.RelationTupleUpdate{
		Operation: core.RelationTupleUpdate_TOUCH,
		Tuple:     tuple.MustParse(relationship),
	}
}

func del(relationship string) *core.RelationTupleUpdate {
	return &core.RelationTupleUpdate{
		Operation: core.RelationTupleUpdate_DELETE,
		Tuple:     tuple.MustParse(relationship),
	}
}

func canonicalize(in []datastore.RevisionChanges) []datastore.RevisionChanges {
	out := make([]datastore.RevisionChanges, 0, len(in))

	for _, rev := range in {
		outChanges := make([]*core.RelationTupleUpdate, 0, len(rev.RelationshipChanges))

		outChanges = append(outChanges, rev.RelationshipChanges...)
		sort.Slice(outChanges, func(i, j int) bool {
			// Return if i < j
			left, right := outChanges[i], outChanges[j]
			tupleCompareResult := strings.Compare(tuple.StringWithoutCaveat(left.Tuple), tuple.StringWithoutCaveat(right.Tuple))
			if tupleCompareResult < 0 {
				return true
			}
			if tupleCompareResult > 0 {
				return false
			}

			// Tuples are equal, sort by op
			return left.Operation < right.Operation
		})

		deletedNamespaces := rev.DeletedNamespaces
		if len(rev.DeletedNamespaces) == 0 {
			deletedNamespaces = nil
		} else {
			sort.Strings(deletedNamespaces)
		}

		deletedCaveats := rev.DeletedCaveats
		if len(rev.DeletedCaveats) == 0 {
			deletedCaveats = nil
		} else {
			sort.Strings(deletedCaveats)
		}

		changedDefinitions := rev.ChangedDefinitions
		if len(rev.ChangedDefinitions) == 0 {
			changedDefinitions = nil
		}

		out = append(out, datastore.RevisionChanges{
			Revision:            rev.Revision,
			RelationshipChanges: outChanges,
			DeletedNamespaces:   deletedNamespaces,
			DeletedCaveats:      deletedCaveats,
			ChangedDefinitions:  changedDefinitions,
		})
	}

	return out
}
