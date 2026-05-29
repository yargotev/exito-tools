package app

import (
	"github.com/yargotev/exito-tools/internal/config"
	"github.com/yargotev/exito-tools/internal/domain/catalog"
	"github.com/yargotev/exito-tools/internal/domain/checkout"
	"github.com/yargotev/exito-tools/internal/domain/geo"
	"github.com/yargotev/exito-tools/internal/domain/orders"
	"github.com/yargotev/exito-tools/internal/registry"
	"github.com/yargotev/exito-tools/internal/workflow"
)

// Options contains application boot inputs shared by all surfaces.
type Options struct {
	Config config.Options
}

// Application is the explicit application wiring seam shared by surfaces.
type Application struct {
	Config        config.Effective
	ConfigOptions config.Options
	Registry      registry.Registry
}

// New builds the minimal application scaffold, resolves configuration, and finalizes the registry.
func New(options Options) (*Application, error) {
	effectiveConfig, err := config.Resolve(options.Config)
	if err != nil {
		return nil, err
	}

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(orders.NewGetCapability(ordersGetter(effectiveConfig))); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(orders.NewGetVTEXOMSCapability(vtexOMSGetter(effectiveConfig))); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(geo.NewGeocodeAddressCapability(geoGeocoder(effectiveConfig))); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(geo.NewResolveVTEXRegionCapability(geoVTEXRegionResolver(effectiveConfig))); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(catalog.NewSearchProductsCapability(catalogSearcher(effectiveConfig))); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(catalog.NewIntelligentSearchProductsCapability(catalogIntelligentSearcher(effectiveConfig))); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(catalog.NewCreateVTEXSegmentCapability(catalogVTEXSegmentCreator(effectiveConfig))); err != nil {
		return nil, err
	}
	checkoutProvider := checkoutClient(effectiveConfig)
	if err := builder.RegisterExecutable(checkout.NewGetOrderFormCapability(checkoutProvider)); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(checkout.NewCreateOrderFormCapability(checkoutProvider)); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(checkout.NewAddItemsCapability(checkoutProvider)); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(checkout.NewUpdateClientProfileCapability(checkoutProvider)); err != nil {
		return nil, err
	}
	if err := builder.RegisterExecutable(workflow.NewRegionalizedIntelligentSearchProductsCapability(
		geoVTEXRegionResolver(effectiveConfig),
		catalogVTEXSegmentCreator(effectiveConfig),
		catalogIntelligentSearcher(effectiveConfig),
	)); err != nil {
		return nil, err
	}

	return &Application{
		Config:        effectiveConfig,
		ConfigOptions: options.Config,
		Registry:      builder.Finalize(),
	}, nil
}

func ordersGetter(effectiveConfig config.Effective) orders.Getter {
	if !effectiveConfig.OrdersProvider.Configured {
		return orders.UnavailableGetter{}
	}

	return orders.NewHTTPGetter(orders.HTTPGetterConfig{
		BaseURL:      effectiveConfig.OrdersProvider.BaseURL,
		Token:        effectiveConfig.OrdersProvider.Token,
		TokenURL:     effectiveConfig.OrdersProvider.TokenURL,
		ClientID:     effectiveConfig.OrdersProvider.ClientID,
		ClientSecret: effectiveConfig.OrdersProvider.ClientSecret,
		Scope:        effectiveConfig.OrdersProvider.Scope,
	}, nil)
}

func vtexOMSGetter(effectiveConfig config.Effective) orders.VTEXOMSGetterPort {
	return orders.NewVTEXOMSBrandGetter(
		vtexOMSBrandGetter(effectiveConfig.VTEXOMSProvider.Exito),
		vtexOMSBrandGetter(effectiveConfig.VTEXOMSProvider.Carulla),
	)
}

func vtexOMSBrandGetter(provider config.VTEXOMSBrandProvider) orders.VTEXOMSGetterPort {
	if !provider.Configured {
		return orders.UnavailableVTEXOMSGetter{}
	}
	return orders.NewVTEXOMSGetter(orders.VTEXOMSGetterConfig{
		BaseURL:  provider.BaseURL,
		AppKey:   provider.AppKey,
		AppToken: provider.AppToken,
	}, nil)
}

func geoGeocoder(effectiveConfig config.Effective) geo.Geocoder {
	if !effectiveConfig.GeoProvider.Configured {
		return geo.UnavailableGeocoder{}
	}

	return geo.NewHTTPGeocoder(geo.HTTPGeocoderConfig{
		BaseURL: effectiveConfig.GeoProvider.BaseURL,
		Token:   effectiveConfig.GeoProvider.Token,
	}, nil)
}

func catalogSearcher(effectiveConfig config.Effective) catalog.Searcher {
	return catalog.NewBrandSearcher(
		catalogBrandSearcher(effectiveConfig.VTEXCatalogProvider.Exito),
		catalogBrandSearcher(effectiveConfig.VTEXCatalogProvider.Carulla),
	)
}

func catalogBrandSearcher(provider config.VTEXCatalogBrandProvider) catalog.Searcher {
	if !provider.Configured {
		return catalog.UnavailableSearcher{}
	}
	return catalog.NewHTTPSearcher(catalog.HTTPSearcherConfig{BaseURL: provider.BaseURL}, nil)
}

func catalogIntelligentSearcher(effectiveConfig config.Effective) catalog.IntelligentSearchProductsSearcher {
	return catalog.NewIntelligentBrandSearcher(
		catalogIntelligentBrandSearcher(effectiveConfig.VTEXIntelligentSearchProvider.Exito),
		catalogIntelligentBrandSearcher(effectiveConfig.VTEXIntelligentSearchProvider.Carulla),
	)
}

func catalogIntelligentBrandSearcher(provider config.VTEXCatalogBrandProvider) catalog.IntelligentSearchProductsSearcher {
	if !provider.Configured {
		return catalog.UnavailableIntelligentSearcher{}
	}
	return catalog.NewHTTPIntelligentSearcher(catalog.HTTPIntelligentSearcherConfig{BaseURL: provider.BaseURL}, nil)
}

func catalogVTEXSegmentCreator(effectiveConfig config.Effective) catalog.VTEXSegmentCreator {
	return catalog.NewVTEXSegmentBrandCreator(
		catalogVTEXBrandSegmentCreator(effectiveConfig.VTEXCatalogProvider.Exito),
		catalogVTEXBrandSegmentCreator(effectiveConfig.VTEXCatalogProvider.Carulla),
	)
}

func catalogVTEXBrandSegmentCreator(provider config.VTEXCatalogBrandProvider) catalog.VTEXSegmentCreator {
	if !provider.Configured {
		return catalog.UnavailableVTEXSegmentCreator{}
	}
	return catalog.NewHTTPVTEXSegmentCreator(catalog.HTTPVTEXSegmentCreatorConfig{BaseURL: provider.BaseURL}, nil)
}

func geoVTEXRegionResolver(effectiveConfig config.Effective) geo.VTEXRegionResolver {
	return geo.NewVTEXBrandRegionResolver(
		geoVTEXBrandRegionResolver(effectiveConfig.VTEXCatalogProvider.Exito),
		geoVTEXBrandRegionResolver(effectiveConfig.VTEXCatalogProvider.Carulla),
	)
}

func geoVTEXBrandRegionResolver(provider config.VTEXCatalogBrandProvider) geo.VTEXRegionResolver {
	if !provider.Configured {
		return geo.UnavailableVTEXRegionResolver{}
	}
	return geo.NewHTTPVTEXRegionResolver(geo.HTTPVTEXRegionResolverConfig{BaseURL: provider.BaseURL}, nil)
}

func checkoutClient(effectiveConfig config.Effective) interface {
	checkout.Getter
	checkout.Creator
	checkout.Adder
	checkout.ClientProfileUpdater
} {
	return checkout.NewBrandClient(
		checkoutBrandClient(effectiveConfig.VTEXCheckoutProvider.Exito),
		checkoutBrandClient(effectiveConfig.VTEXCheckoutProvider.Carulla),
	)
}

func checkoutBrandClient(provider config.VTEXCatalogBrandProvider) interface {
	checkout.Getter
	checkout.Creator
	checkout.Adder
	checkout.ClientProfileUpdater
} {
	if !provider.Configured {
		return checkout.UnavailableClient{}
	}
	return checkout.NewHTTPClient(checkout.HTTPClientConfig{BaseURL: provider.BaseURL}, nil)
}
