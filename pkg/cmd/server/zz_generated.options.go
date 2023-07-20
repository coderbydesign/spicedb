// Code generated by github.com/ecordell/optgen. DO NOT EDIT.
package server

import (
	dispatch "github.com/authzed/spicedb/internal/dispatch"
	graph "github.com/authzed/spicedb/internal/dispatch/graph"
	datastore "github.com/authzed/spicedb/pkg/cmd/datastore"
	util "github.com/authzed/spicedb/pkg/cmd/util"
	datastore1 "github.com/authzed/spicedb/pkg/datastore"
	defaults "github.com/creasty/defaults"
	helpers "github.com/ecordell/optgen/helpers"
	auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	grpc "google.golang.org/grpc"
	"time"
)

type ConfigOption func(c *Config)

// NewConfigWithOptions creates a new Config with the passed in options set
func NewConfigWithOptions(opts ...ConfigOption) *Config {
	c := &Config{}
	for _, o := range opts {
		o(c)
	}
	return c
}

// NewConfigWithOptionsAndDefaults creates a new Config with the passed in options set starting from the defaults
func NewConfigWithOptionsAndDefaults(opts ...ConfigOption) *Config {
	c := &Config{}
	defaults.MustSet(c)
	for _, o := range opts {
		o(c)
	}
	return c
}

// ToOption returns a new ConfigOption that sets the values from the passed in Config
func (c *Config) ToOption() ConfigOption {
	return func(to *Config) {
		to.GRPCServer = c.GRPCServer
		to.GRPCAuthFunc = c.GRPCAuthFunc
		to.PresharedSecureKey = c.PresharedSecureKey
		to.ShutdownGracePeriod = c.ShutdownGracePeriod
		to.DisableVersionResponse = c.DisableVersionResponse
		to.HTTPGateway = c.HTTPGateway
		to.HTTPGatewayUpstreamAddr = c.HTTPGatewayUpstreamAddr
		to.HTTPGatewayUpstreamTLSCertPath = c.HTTPGatewayUpstreamTLSCertPath
		to.HTTPGatewayCorsEnabled = c.HTTPGatewayCorsEnabled
		to.HTTPGatewayCorsAllowedOrigins = c.HTTPGatewayCorsAllowedOrigins
		to.DatastoreConfig = c.DatastoreConfig
		to.Datastore = c.Datastore
		to.MaxCaveatContextSize = c.MaxCaveatContextSize
		to.MaxRelationshipContextSize = c.MaxRelationshipContextSize
		to.NamespaceCacheConfig = c.NamespaceCacheConfig
		to.SchemaPrefixesRequired = c.SchemaPrefixesRequired
		to.DispatchServer = c.DispatchServer
		to.DispatchMaxDepth = c.DispatchMaxDepth
		to.GlobalDispatchConcurrencyLimit = c.GlobalDispatchConcurrencyLimit
		to.DispatchConcurrencyLimits = c.DispatchConcurrencyLimits
		to.DispatchUpstreamAddr = c.DispatchUpstreamAddr
		to.DispatchUpstreamCAPath = c.DispatchUpstreamCAPath
		to.DispatchUpstreamTimeout = c.DispatchUpstreamTimeout
		to.DispatchClientMetricsEnabled = c.DispatchClientMetricsEnabled
		to.DispatchClientMetricsPrefix = c.DispatchClientMetricsPrefix
		to.DispatchClusterMetricsEnabled = c.DispatchClusterMetricsEnabled
		to.DispatchClusterMetricsPrefix = c.DispatchClusterMetricsPrefix
		to.Dispatcher = c.Dispatcher
		to.DispatchHashringReplicationFactor = c.DispatchHashringReplicationFactor
		to.DispatchHashringSpread = c.DispatchHashringSpread
		to.DispatchCacheConfig = c.DispatchCacheConfig
		to.ClusterDispatchCacheConfig = c.ClusterDispatchCacheConfig
		to.DisableV1SchemaAPI = c.DisableV1SchemaAPI
		to.V1SchemaAdditiveOnly = c.V1SchemaAdditiveOnly
		to.MaximumUpdatesPerWrite = c.MaximumUpdatesPerWrite
		to.MaximumPreconditionCount = c.MaximumPreconditionCount
		to.MaxDatastoreReadPageSize = c.MaxDatastoreReadPageSize
		to.StreamingAPITimeout = c.StreamingAPITimeout
		to.MetricsAPI = c.MetricsAPI
		to.UnaryMiddlewareModification = c.UnaryMiddlewareModification
		to.StreamingMiddlewareModification = c.StreamingMiddlewareModification
		to.DispatchUnaryMiddleware = c.DispatchUnaryMiddleware
		to.DispatchStreamingMiddleware = c.DispatchStreamingMiddleware
		to.SilentlyDisableTelemetry = c.SilentlyDisableTelemetry
		to.TelemetryCAOverridePath = c.TelemetryCAOverridePath
		to.TelemetryEndpoint = c.TelemetryEndpoint
		to.TelemetryInterval = c.TelemetryInterval
	}
}

// DebugMap returns a map form of Config for debugging
func (c Config) DebugMap() map[string]any {
	debugMap := map[string]any{}
	debugMap["GRPCServer"] = helpers.DebugValue(c.GRPCServer, false)
	debugMap["GRPCAuthFunc"] = helpers.DebugValue(c.GRPCAuthFunc, false)
	debugMap["PresharedSecureKey"] = helpers.SensitiveDebugValue(c.PresharedSecureKey)
	debugMap["ShutdownGracePeriod"] = helpers.DebugValue(c.ShutdownGracePeriod, false)
	debugMap["DisableVersionResponse"] = helpers.DebugValue(c.DisableVersionResponse, false)
	debugMap["HTTPGateway"] = helpers.DebugValue(c.HTTPGateway, false)
	debugMap["HTTPGatewayUpstreamAddr"] = helpers.DebugValue(c.HTTPGatewayUpstreamAddr, false)
	debugMap["HTTPGatewayUpstreamTLSCertPath"] = helpers.DebugValue(c.HTTPGatewayUpstreamTLSCertPath, false)
	debugMap["HTTPGatewayCorsEnabled"] = helpers.DebugValue(c.HTTPGatewayCorsEnabled, false)
	debugMap["HTTPGatewayCorsAllowedOrigins"] = helpers.DebugValue(c.HTTPGatewayCorsAllowedOrigins, true)
	debugMap["DatastoreConfig"] = helpers.DebugValue(c.DatastoreConfig, false)
	debugMap["Datastore"] = helpers.DebugValue(c.Datastore, false)
	debugMap["MaxCaveatContextSize"] = helpers.DebugValue(c.MaxCaveatContextSize, false)
	debugMap["MaxRelationshipContextSize"] = helpers.DebugValue(c.MaxRelationshipContextSize, false)
	debugMap["NamespaceCacheConfig"] = helpers.DebugValue(c.NamespaceCacheConfig, false)
	debugMap["SchemaPrefixesRequired"] = helpers.DebugValue(c.SchemaPrefixesRequired, false)
	debugMap["DispatchServer"] = helpers.DebugValue(c.DispatchServer, false)
	debugMap["DispatchMaxDepth"] = helpers.DebugValue(c.DispatchMaxDepth, false)
	debugMap["GlobalDispatchConcurrencyLimit"] = helpers.DebugValue(c.GlobalDispatchConcurrencyLimit, false)
	debugMap["DispatchConcurrencyLimits"] = helpers.DebugValue(c.DispatchConcurrencyLimits, false)
	debugMap["DispatchUpstreamAddr"] = helpers.DebugValue(c.DispatchUpstreamAddr, false)
	debugMap["DispatchUpstreamCAPath"] = helpers.DebugValue(c.DispatchUpstreamCAPath, false)
	debugMap["DispatchUpstreamTimeout"] = helpers.DebugValue(c.DispatchUpstreamTimeout, false)
	debugMap["DispatchClientMetricsEnabled"] = helpers.DebugValue(c.DispatchClientMetricsEnabled, false)
	debugMap["DispatchClientMetricsPrefix"] = helpers.DebugValue(c.DispatchClientMetricsPrefix, false)
	debugMap["DispatchClusterMetricsEnabled"] = helpers.DebugValue(c.DispatchClusterMetricsEnabled, false)
	debugMap["DispatchClusterMetricsPrefix"] = helpers.DebugValue(c.DispatchClusterMetricsPrefix, false)
	debugMap["Dispatcher"] = helpers.DebugValue(c.Dispatcher, false)
	debugMap["DispatchHashringReplicationFactor"] = helpers.DebugValue(c.DispatchHashringReplicationFactor, false)
	debugMap["DispatchHashringSpread"] = helpers.DebugValue(c.DispatchHashringSpread, false)
	debugMap["DispatchCacheConfig"] = helpers.DebugValue(c.DispatchCacheConfig, false)
	debugMap["ClusterDispatchCacheConfig"] = helpers.DebugValue(c.ClusterDispatchCacheConfig, false)
	debugMap["DisableV1SchemaAPI"] = helpers.DebugValue(c.DisableV1SchemaAPI, false)
	debugMap["V1SchemaAdditiveOnly"] = helpers.DebugValue(c.V1SchemaAdditiveOnly, false)
	debugMap["MaximumUpdatesPerWrite"] = helpers.DebugValue(c.MaximumUpdatesPerWrite, false)
	debugMap["MaximumPreconditionCount"] = helpers.DebugValue(c.MaximumPreconditionCount, false)
	debugMap["MaxDatastoreReadPageSize"] = helpers.DebugValue(c.MaxDatastoreReadPageSize, false)
	debugMap["StreamingAPITimeout"] = helpers.DebugValue(c.StreamingAPITimeout, false)
	debugMap["MetricsAPI"] = helpers.DebugValue(c.MetricsAPI, false)
	debugMap["SilentlyDisableTelemetry"] = helpers.DebugValue(c.SilentlyDisableTelemetry, false)
	debugMap["TelemetryCAOverridePath"] = helpers.DebugValue(c.TelemetryCAOverridePath, false)
	debugMap["TelemetryEndpoint"] = helpers.DebugValue(c.TelemetryEndpoint, false)
	debugMap["TelemetryInterval"] = helpers.DebugValue(c.TelemetryInterval, false)
	return debugMap
}

// ConfigWithOptions configures an existing Config with the passed in options set
func ConfigWithOptions(c *Config, opts ...ConfigOption) *Config {
	for _, o := range opts {
		o(c)
	}
	return c
}

// WithOptions configures the receiver Config with the passed in options set
func (c *Config) WithOptions(opts ...ConfigOption) *Config {
	for _, o := range opts {
		o(c)
	}
	return c
}

// WithGRPCServer returns an option that can set GRPCServer on a Config
func WithGRPCServer(gRPCServer util.GRPCServerConfig) ConfigOption {
	return func(c *Config) {
		c.GRPCServer = gRPCServer
	}
}

// WithGRPCAuthFunc returns an option that can set GRPCAuthFunc on a Config
func WithGRPCAuthFunc(gRPCAuthFunc auth.AuthFunc) ConfigOption {
	return func(c *Config) {
		c.GRPCAuthFunc = gRPCAuthFunc
	}
}

// WithPresharedSecureKey returns an option that can append PresharedSecureKeys to Config.PresharedSecureKey
func WithPresharedSecureKey(presharedSecureKey string) ConfigOption {
	return func(c *Config) {
		c.PresharedSecureKey = append(c.PresharedSecureKey, presharedSecureKey)
	}
}

// SetPresharedSecureKey returns an option that can set PresharedSecureKey on a Config
func SetPresharedSecureKey(presharedSecureKey []string) ConfigOption {
	return func(c *Config) {
		c.PresharedSecureKey = presharedSecureKey
	}
}

// WithShutdownGracePeriod returns an option that can set ShutdownGracePeriod on a Config
func WithShutdownGracePeriod(shutdownGracePeriod time.Duration) ConfigOption {
	return func(c *Config) {
		c.ShutdownGracePeriod = shutdownGracePeriod
	}
}

// WithDisableVersionResponse returns an option that can set DisableVersionResponse on a Config
func WithDisableVersionResponse(disableVersionResponse bool) ConfigOption {
	return func(c *Config) {
		c.DisableVersionResponse = disableVersionResponse
	}
}

// WithHTTPGateway returns an option that can set HTTPGateway on a Config
func WithHTTPGateway(hTTPGateway util.HTTPServerConfig) ConfigOption {
	return func(c *Config) {
		c.HTTPGateway = hTTPGateway
	}
}

// WithHTTPGatewayUpstreamAddr returns an option that can set HTTPGatewayUpstreamAddr on a Config
func WithHTTPGatewayUpstreamAddr(hTTPGatewayUpstreamAddr string) ConfigOption {
	return func(c *Config) {
		c.HTTPGatewayUpstreamAddr = hTTPGatewayUpstreamAddr
	}
}

// WithHTTPGatewayUpstreamTLSCertPath returns an option that can set HTTPGatewayUpstreamTLSCertPath on a Config
func WithHTTPGatewayUpstreamTLSCertPath(hTTPGatewayUpstreamTLSCertPath string) ConfigOption {
	return func(c *Config) {
		c.HTTPGatewayUpstreamTLSCertPath = hTTPGatewayUpstreamTLSCertPath
	}
}

// WithHTTPGatewayCorsEnabled returns an option that can set HTTPGatewayCorsEnabled on a Config
func WithHTTPGatewayCorsEnabled(hTTPGatewayCorsEnabled bool) ConfigOption {
	return func(c *Config) {
		c.HTTPGatewayCorsEnabled = hTTPGatewayCorsEnabled
	}
}

// WithHTTPGatewayCorsAllowedOrigins returns an option that can append HTTPGatewayCorsAllowedOriginss to Config.HTTPGatewayCorsAllowedOrigins
func WithHTTPGatewayCorsAllowedOrigins(hTTPGatewayCorsAllowedOrigins string) ConfigOption {
	return func(c *Config) {
		c.HTTPGatewayCorsAllowedOrigins = append(c.HTTPGatewayCorsAllowedOrigins, hTTPGatewayCorsAllowedOrigins)
	}
}

// SetHTTPGatewayCorsAllowedOrigins returns an option that can set HTTPGatewayCorsAllowedOrigins on a Config
func SetHTTPGatewayCorsAllowedOrigins(hTTPGatewayCorsAllowedOrigins []string) ConfigOption {
	return func(c *Config) {
		c.HTTPGatewayCorsAllowedOrigins = hTTPGatewayCorsAllowedOrigins
	}
}

// WithDatastoreConfig returns an option that can set DatastoreConfig on a Config
func WithDatastoreConfig(datastoreConfig datastore.Config) ConfigOption {
	return func(c *Config) {
		c.DatastoreConfig = datastoreConfig
	}
}

// WithDatastore returns an option that can set Datastore on a Config
func WithDatastore(datastore datastore1.Datastore) ConfigOption {
	return func(c *Config) {
		c.Datastore = datastore
	}
}

// WithMaxCaveatContextSize returns an option that can set MaxCaveatContextSize on a Config
func WithMaxCaveatContextSize(maxCaveatContextSize int) ConfigOption {
	return func(c *Config) {
		c.MaxCaveatContextSize = maxCaveatContextSize
	}
}

// WithMaxRelationshipContextSize returns an option that can set MaxRelationshipContextSize on a Config
func WithMaxRelationshipContextSize(maxRelationshipContextSize int) ConfigOption {
	return func(c *Config) {
		c.MaxRelationshipContextSize = maxRelationshipContextSize
	}
}

// WithNamespaceCacheConfig returns an option that can set NamespaceCacheConfig on a Config
func WithNamespaceCacheConfig(namespaceCacheConfig CacheConfig) ConfigOption {
	return func(c *Config) {
		c.NamespaceCacheConfig = namespaceCacheConfig
	}
}

// WithSchemaPrefixesRequired returns an option that can set SchemaPrefixesRequired on a Config
func WithSchemaPrefixesRequired(schemaPrefixesRequired bool) ConfigOption {
	return func(c *Config) {
		c.SchemaPrefixesRequired = schemaPrefixesRequired
	}
}

// WithDispatchServer returns an option that can set DispatchServer on a Config
func WithDispatchServer(dispatchServer util.GRPCServerConfig) ConfigOption {
	return func(c *Config) {
		c.DispatchServer = dispatchServer
	}
}

// WithDispatchMaxDepth returns an option that can set DispatchMaxDepth on a Config
func WithDispatchMaxDepth(dispatchMaxDepth uint32) ConfigOption {
	return func(c *Config) {
		c.DispatchMaxDepth = dispatchMaxDepth
	}
}

// WithGlobalDispatchConcurrencyLimit returns an option that can set GlobalDispatchConcurrencyLimit on a Config
func WithGlobalDispatchConcurrencyLimit(globalDispatchConcurrencyLimit uint16) ConfigOption {
	return func(c *Config) {
		c.GlobalDispatchConcurrencyLimit = globalDispatchConcurrencyLimit
	}
}

// WithDispatchConcurrencyLimits returns an option that can set DispatchConcurrencyLimits on a Config
func WithDispatchConcurrencyLimits(dispatchConcurrencyLimits graph.ConcurrencyLimits) ConfigOption {
	return func(c *Config) {
		c.DispatchConcurrencyLimits = dispatchConcurrencyLimits
	}
}

// WithDispatchUpstreamAddr returns an option that can set DispatchUpstreamAddr on a Config
func WithDispatchUpstreamAddr(dispatchUpstreamAddr string) ConfigOption {
	return func(c *Config) {
		c.DispatchUpstreamAddr = dispatchUpstreamAddr
	}
}

// WithDispatchUpstreamCAPath returns an option that can set DispatchUpstreamCAPath on a Config
func WithDispatchUpstreamCAPath(dispatchUpstreamCAPath string) ConfigOption {
	return func(c *Config) {
		c.DispatchUpstreamCAPath = dispatchUpstreamCAPath
	}
}

// WithDispatchUpstreamTimeout returns an option that can set DispatchUpstreamTimeout on a Config
func WithDispatchUpstreamTimeout(dispatchUpstreamTimeout time.Duration) ConfigOption {
	return func(c *Config) {
		c.DispatchUpstreamTimeout = dispatchUpstreamTimeout
	}
}

// WithDispatchClientMetricsEnabled returns an option that can set DispatchClientMetricsEnabled on a Config
func WithDispatchClientMetricsEnabled(dispatchClientMetricsEnabled bool) ConfigOption {
	return func(c *Config) {
		c.DispatchClientMetricsEnabled = dispatchClientMetricsEnabled
	}
}

// WithDispatchClientMetricsPrefix returns an option that can set DispatchClientMetricsPrefix on a Config
func WithDispatchClientMetricsPrefix(dispatchClientMetricsPrefix string) ConfigOption {
	return func(c *Config) {
		c.DispatchClientMetricsPrefix = dispatchClientMetricsPrefix
	}
}

// WithDispatchClusterMetricsEnabled returns an option that can set DispatchClusterMetricsEnabled on a Config
func WithDispatchClusterMetricsEnabled(dispatchClusterMetricsEnabled bool) ConfigOption {
	return func(c *Config) {
		c.DispatchClusterMetricsEnabled = dispatchClusterMetricsEnabled
	}
}

// WithDispatchClusterMetricsPrefix returns an option that can set DispatchClusterMetricsPrefix on a Config
func WithDispatchClusterMetricsPrefix(dispatchClusterMetricsPrefix string) ConfigOption {
	return func(c *Config) {
		c.DispatchClusterMetricsPrefix = dispatchClusterMetricsPrefix
	}
}

// WithDispatcher returns an option that can set Dispatcher on a Config
func WithDispatcher(dispatcher dispatch.Dispatcher) ConfigOption {
	return func(c *Config) {
		c.Dispatcher = dispatcher
	}
}

// WithDispatchHashringReplicationFactor returns an option that can set DispatchHashringReplicationFactor on a Config
func WithDispatchHashringReplicationFactor(dispatchHashringReplicationFactor uint16) ConfigOption {
	return func(c *Config) {
		c.DispatchHashringReplicationFactor = dispatchHashringReplicationFactor
	}
}

// WithDispatchHashringSpread returns an option that can set DispatchHashringSpread on a Config
func WithDispatchHashringSpread(dispatchHashringSpread uint8) ConfigOption {
	return func(c *Config) {
		c.DispatchHashringSpread = dispatchHashringSpread
	}
}

// WithDispatchCacheConfig returns an option that can set DispatchCacheConfig on a Config
func WithDispatchCacheConfig(dispatchCacheConfig CacheConfig) ConfigOption {
	return func(c *Config) {
		c.DispatchCacheConfig = dispatchCacheConfig
	}
}

// WithClusterDispatchCacheConfig returns an option that can set ClusterDispatchCacheConfig on a Config
func WithClusterDispatchCacheConfig(clusterDispatchCacheConfig CacheConfig) ConfigOption {
	return func(c *Config) {
		c.ClusterDispatchCacheConfig = clusterDispatchCacheConfig
	}
}

// WithDisableV1SchemaAPI returns an option that can set DisableV1SchemaAPI on a Config
func WithDisableV1SchemaAPI(disableV1SchemaAPI bool) ConfigOption {
	return func(c *Config) {
		c.DisableV1SchemaAPI = disableV1SchemaAPI
	}
}

// WithV1SchemaAdditiveOnly returns an option that can set V1SchemaAdditiveOnly on a Config
func WithV1SchemaAdditiveOnly(v1SchemaAdditiveOnly bool) ConfigOption {
	return func(c *Config) {
		c.V1SchemaAdditiveOnly = v1SchemaAdditiveOnly
	}
}

// WithMaximumUpdatesPerWrite returns an option that can set MaximumUpdatesPerWrite on a Config
func WithMaximumUpdatesPerWrite(maximumUpdatesPerWrite uint16) ConfigOption {
	return func(c *Config) {
		c.MaximumUpdatesPerWrite = maximumUpdatesPerWrite
	}
}

// WithMaximumPreconditionCount returns an option that can set MaximumPreconditionCount on a Config
func WithMaximumPreconditionCount(maximumPreconditionCount uint16) ConfigOption {
	return func(c *Config) {
		c.MaximumPreconditionCount = maximumPreconditionCount
	}
}

// WithMaxDatastoreReadPageSize returns an option that can set MaxDatastoreReadPageSize on a Config
func WithMaxDatastoreReadPageSize(maxDatastoreReadPageSize uint64) ConfigOption {
	return func(c *Config) {
		c.MaxDatastoreReadPageSize = maxDatastoreReadPageSize
	}
}

// WithStreamingAPITimeout returns an option that can set StreamingAPITimeout on a Config
func WithStreamingAPITimeout(streamingAPITimeout time.Duration) ConfigOption {
	return func(c *Config) {
		c.StreamingAPITimeout = streamingAPITimeout
	}
}

// WithMetricsAPI returns an option that can set MetricsAPI on a Config
func WithMetricsAPI(metricsAPI util.HTTPServerConfig) ConfigOption {
	return func(c *Config) {
		c.MetricsAPI = metricsAPI
	}
}

// WithUnaryMiddlewareModification returns an option that can append UnaryMiddlewareModifications to Config.UnaryMiddlewareModification
func WithUnaryMiddlewareModification(unaryMiddlewareModification MiddlewareModification[grpc.UnaryServerInterceptor]) ConfigOption {
	return func(c *Config) {
		c.UnaryMiddlewareModification = append(c.UnaryMiddlewareModification, unaryMiddlewareModification)
	}
}

// SetUnaryMiddlewareModification returns an option that can set UnaryMiddlewareModification on a Config
func SetUnaryMiddlewareModification(unaryMiddlewareModification []MiddlewareModification[grpc.UnaryServerInterceptor]) ConfigOption {
	return func(c *Config) {
		c.UnaryMiddlewareModification = unaryMiddlewareModification
	}
}

// WithStreamingMiddlewareModification returns an option that can append StreamingMiddlewareModifications to Config.StreamingMiddlewareModification
func WithStreamingMiddlewareModification(streamingMiddlewareModification MiddlewareModification[grpc.StreamServerInterceptor]) ConfigOption {
	return func(c *Config) {
		c.StreamingMiddlewareModification = append(c.StreamingMiddlewareModification, streamingMiddlewareModification)
	}
}

// SetStreamingMiddlewareModification returns an option that can set StreamingMiddlewareModification on a Config
func SetStreamingMiddlewareModification(streamingMiddlewareModification []MiddlewareModification[grpc.StreamServerInterceptor]) ConfigOption {
	return func(c *Config) {
		c.StreamingMiddlewareModification = streamingMiddlewareModification
	}
}

// WithDispatchUnaryMiddleware returns an option that can append DispatchUnaryMiddlewares to Config.DispatchUnaryMiddleware
func WithDispatchUnaryMiddleware(dispatchUnaryMiddleware grpc.UnaryServerInterceptor) ConfigOption {
	return func(c *Config) {
		c.DispatchUnaryMiddleware = append(c.DispatchUnaryMiddleware, dispatchUnaryMiddleware)
	}
}

// SetDispatchUnaryMiddleware returns an option that can set DispatchUnaryMiddleware on a Config
func SetDispatchUnaryMiddleware(dispatchUnaryMiddleware []grpc.UnaryServerInterceptor) ConfigOption {
	return func(c *Config) {
		c.DispatchUnaryMiddleware = dispatchUnaryMiddleware
	}
}

// WithDispatchStreamingMiddleware returns an option that can append DispatchStreamingMiddlewares to Config.DispatchStreamingMiddleware
func WithDispatchStreamingMiddleware(dispatchStreamingMiddleware grpc.StreamServerInterceptor) ConfigOption {
	return func(c *Config) {
		c.DispatchStreamingMiddleware = append(c.DispatchStreamingMiddleware, dispatchStreamingMiddleware)
	}
}

// SetDispatchStreamingMiddleware returns an option that can set DispatchStreamingMiddleware on a Config
func SetDispatchStreamingMiddleware(dispatchStreamingMiddleware []grpc.StreamServerInterceptor) ConfigOption {
	return func(c *Config) {
		c.DispatchStreamingMiddleware = dispatchStreamingMiddleware
	}
}

// WithSilentlyDisableTelemetry returns an option that can set SilentlyDisableTelemetry on a Config
func WithSilentlyDisableTelemetry(silentlyDisableTelemetry bool) ConfigOption {
	return func(c *Config) {
		c.SilentlyDisableTelemetry = silentlyDisableTelemetry
	}
}

// WithTelemetryCAOverridePath returns an option that can set TelemetryCAOverridePath on a Config
func WithTelemetryCAOverridePath(telemetryCAOverridePath string) ConfigOption {
	return func(c *Config) {
		c.TelemetryCAOverridePath = telemetryCAOverridePath
	}
}

// WithTelemetryEndpoint returns an option that can set TelemetryEndpoint on a Config
func WithTelemetryEndpoint(telemetryEndpoint string) ConfigOption {
	return func(c *Config) {
		c.TelemetryEndpoint = telemetryEndpoint
	}
}

// WithTelemetryInterval returns an option that can set TelemetryInterval on a Config
func WithTelemetryInterval(telemetryInterval time.Duration) ConfigOption {
	return func(c *Config) {
		c.TelemetryInterval = telemetryInterval
	}
}
