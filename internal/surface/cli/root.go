package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/config"
	"github.com/yargotev/exito-tools/internal/domain/catalog"
	"github.com/yargotev/exito-tools/internal/domain/checkout"
	"github.com/yargotev/exito-tools/internal/domain/geo"
	"github.com/yargotev/exito-tools/internal/domain/orders"
	"github.com/yargotev/exito-tools/internal/execution"
	"github.com/yargotev/exito-tools/internal/presenter"
	tuisurface "github.com/yargotev/exito-tools/internal/surface/tui"
	"github.com/yargotev/exito-tools/internal/workflow"
)

// Bootstrapper builds the application after Cobra has parsed CLI boot flags.
type Bootstrapper func(app.Options) (*app.Application, error)

type rootOptions struct {
	configPath    string
	profile       string
	correlationID string
}

type capabilitiesData struct {
	Capabilities []capability.Definition `json:"capabilities"`
}

type defaultProfileData struct {
	Profile      string        `json:"profile"`
	ConfigPath   string        `json:"configPath"`
	ConfigSource config.Source `json:"configSource"`
}

// NewRoot builds the minimal English-only CLI root surface.
func NewRoot(bootstrap Bootstrapper) *cobra.Command {
	if bootstrap == nil {
		bootstrap = app.New
	}

	options := rootOptions{}
	command := &cobra.Command{
		Use:   "exito",
		Short: "Exito Tools command-line interface",
		Long:  rootLong(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(options))
			if err != nil {
				return err
			}

			cmd.Long = rootLong(len(application.Registry.All()))
			return cmd.Help()
		},
		SilenceErrors:     true,
		SilenceUsage:      true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	command.PersistentFlags().StringVar(&options.configPath, "config", "", "Path to the Exito Tools configuration file")
	command.PersistentFlags().StringVar(&options.profile, "profile", "", "Configuration profile to use")
	command.PersistentFlags().StringVar(&options.correlationID, "correlation-id", "", "Correlation ID to include in JSON command metadata")
	command.SetHelpCommand(&cobra.Command{Hidden: true})
	command.AddCommand(newCapabilitiesCommand(bootstrap, &options))
	command.AddCommand(newConfigCommand(&options))
	command.AddCommand(newRunCommand(bootstrap, &options))
	command.AddCommand(newOrdersCommand(bootstrap, &options))
	command.AddCommand(newGeoCommand(bootstrap, &options))
	command.AddCommand(newCatalogCommand(bootstrap, &options))
	command.AddCommand(newCheckoutCommand(bootstrap, &options))
	command.AddCommand(newTUICommand(bootstrap, &options))
	return command
}

func newConfigCommand(options *rootOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "Manage non-sensitive application configuration",
	}
	command.AddCommand(newSetDefaultProfileCommand(options))
	return command
}

func newSetDefaultProfileCommand(options *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "set-default-profile <profile>",
		Short: "Persist the saved Default Profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			startedAt := time.Now()
			requestID, err := execution.NewRequestID()
			if err != nil {
				return err
			}

			result, err := config.SetDefaultProfile(appOptions(*options).Config, args[0])
			if err != nil {
				return err
			}

			data := defaultProfileData{
				Profile:      result.Profile,
				ConfigPath:   result.ConfigPath,
				ConfigSource: result.ConfigSource,
			}
			metadata := execution.NewMetadata(requestID, options.correlationID, startedAt, time.Now())
			envelope := capability.Envelope[defaultProfileData]{
				OK:   true,
				Data: &data,
				Meta: metadata.EnvelopeMeta(result.Profile, ""),
			}

			return presenter.WriteJSON(cmd.OutOrStdout(), envelope)
		},
	}
}

func newTUICommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Open the interactive TUI",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			return tuisurface.Run(cmd.Context(), application, tuisurface.IO{
				Input:  cmd.InOrStdin(),
				Output: cmd.OutOrStdout(),
			})
		},
	}
}

func newCapabilitiesCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "capabilities",
		Short: "Print the machine-readable capability inventory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			startedAt := time.Now()
			requestID, err := execution.NewRequestID()
			if err != nil {
				return err
			}

			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			data := capabilitiesData{Capabilities: application.Registry.All()}
			metadata := execution.NewMetadata(requestID, options.correlationID, startedAt, time.Now())
			envelope := capability.Envelope[capabilitiesData]{
				OK:   true,
				Data: &data,
				Meta: metadata.EnvelopeMeta(application.Config.Profile, ""),
			}

			return presenter.WriteJSON(cmd.OutOrStdout(), envelope)
		},
	}
}

func newCheckoutCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "checkout",
		Short: "Run Checkout Domain commands",
	}
	command.AddCommand(newCheckoutGetOrderFormCommand(bootstrap, options))
	command.AddCommand(newCheckoutCreateOrderFormCommand(bootstrap, options))
	command.AddCommand(newCheckoutAddItemsCommand(bootstrap, options))
	return command
}

func newCheckoutGetOrderFormCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var brand string
	var orderFormID string

	command := &cobra.Command{
		Use:   "get-order-form",
		Short: "Get a VTEX Checkout orderForm by ID",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID: checkout.CapabilityGetOrderFormID,
				Input: capability.Input{
					"brand":       brand,
					"orderFormId": orderFormID,
				},
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&brand, "brand", "exito", "VTEX brand account to query: exito or carulla")
	command.Flags().StringVar(&orderFormID, "order-form-id", "", "VTEX Checkout orderForm identifier")
	_ = command.MarkFlagRequired("order-form-id")
	return command
}

func newCheckoutCreateOrderFormCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var brand string
	var salesChannel string
	var confirmed bool

	command := &cobra.Command{
		Use:   "create-order-form",
		Short: "Create a VTEX Checkout orderForm",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID: checkout.CapabilityCreateOrderFormID,
				Input: capability.Input{
					"brand":        brand,
					"salesChannel": salesChannel,
				},
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
				Confirmed:     confirmed,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&brand, "brand", "exito", "VTEX brand account to query: exito or carulla")
	command.Flags().StringVar(&salesChannel, "sales-channel", "", "VTEX sales channel/trade policy used as sc")
	command.Flags().BoolVar(&confirmed, "confirm", false, "Explicitly confirm VTEX Checkout orderForm creation")
	_ = command.MarkFlagRequired("sales-channel")
	return command
}

func newCheckoutAddItemsCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var brand string
	var orderFormID string
	var rawItems []string
	var confirmed bool

	command := &cobra.Command{
		Use:   "add-items",
		Short: "Add SKU items to a VTEX Checkout orderForm",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := parseCheckoutItems(rawItems)
			if err != nil {
				return err
			}

			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID: checkout.CapabilityAddItemsID,
				Input: capability.Input{
					"brand":       brand,
					"orderFormId": orderFormID,
					"items":       checkoutItemsAsCapabilityInput(items),
				},
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
				Confirmed:     confirmed,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&brand, "brand", "exito", "VTEX brand account to query: exito or carulla")
	command.Flags().StringVar(&orderFormID, "order-form-id", "", "VTEX Checkout orderForm identifier")
	command.Flags().StringArrayVar(&rawItems, "item", nil, "Item to add as sku=<sku>,quantity=<qty>[,seller=<seller>]; repeat for multiple items")
	command.Flags().BoolVar(&confirmed, "confirm", false, "Explicitly confirm VTEX Checkout item addition")
	_ = command.MarkFlagRequired("order-form-id")
	_ = command.MarkFlagRequired("item")
	return command
}

func parseCheckoutItems(rawItems []string) ([]checkout.AddItemInput, error) {
	items := make([]checkout.AddItemInput, 0, len(rawItems))
	for _, raw := range rawItems {
		fields := map[string]string{}
		for _, part := range strings.Split(raw, ",") {
			key, value, ok := strings.Cut(part, "=")
			if !ok {
				return nil, fmt.Errorf("invalid --item %q: expected sku=<sku>,quantity=<qty>[,seller=<seller>]", raw)
			}
			fields[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
		quantity, err := strconv.Atoi(fields["quantity"])
		if err != nil {
			return nil, fmt.Errorf("invalid --item %q: quantity must be an integer", raw)
		}
		seller := fields["seller"]
		if seller == "" {
			seller = "1"
		}
		items = append(items, checkout.AddItemInput{SKU: fields["sku"], Quantity: quantity, Seller: seller})
	}
	return items, nil
}

func checkoutItemsAsCapabilityInput(items []checkout.AddItemInput) []any {
	out := make([]any, 0, len(items))
	for _, item := range items {
		out = append(out, map[string]any{"sku": item.SKU, "quantity": item.Quantity, "seller": item.Seller})
	}
	return out
}

func newCatalogCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "catalog",
		Short: "Run Catalog Domain commands",
	}
	command.AddCommand(newCatalogSearchProductsCommand(bootstrap, options))
	command.AddCommand(newCatalogIntelligentSearchCommand(bootstrap, options))
	command.AddCommand(newCatalogCreateVTEXSegmentCommand(bootstrap, options))
	return command
}

func newCatalogSearchProductsCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var brand string
	var by string
	var value string
	var filters []string
	var fullText string
	var order string
	var from int
	var to int

	command := &cobra.Command{
		Use:   "search-products",
		Short: "Search VTEX catalog products",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			input := capability.Input{
				"brand": brand,
				"from":  from,
				"to":    to,
			}
			if by != "" {
				input["by"] = by
			}
			if value != "" {
				input["value"] = value
			}
			if len(filters) > 0 {
				input["fq"] = filters
			}
			if fullText != "" {
				input["ft"] = fullText
			}
			if order != "" {
				input["order"] = order
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID:  catalog.CapabilitySearchProductsID,
				Input:         input,
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&brand, "brand", "exito", "VTEX brand account to query: exito or carulla")
	command.Flags().StringVar(&by, "by", "", "Friendly lookup mode: sku-id, product-id, ref-id, ean, seller-id, category, brand-id, collection-id, text, or slug")
	command.Flags().StringVar(&value, "value", "", "Lookup value used with --by")
	command.Flags().StringArrayVar(&filters, "fq", nil, "Raw VTEX fq filter; repeat for multiple filters")
	command.Flags().StringVar(&fullText, "ft", "", "VTEX full-text search term")
	command.Flags().StringVar(&order, "order", "", "VTEX O sorting value, such as OrderByPriceASC")
	command.Flags().IntVar(&from, "from", 0, "Initial result index")
	command.Flags().IntVar(&to, "to", 9, "Final result index; VTEX allows at most 50 results per request")
	return command
}

func newCatalogCreateVTEXSegmentCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var brand string
	var regionID string
	var salesChannel string
	var includeCookie bool
	var confirmed bool

	command := &cobra.Command{
		Use:   "create-vtex-segment",
		Short: "Create a VTEX segment token from a region ID and sales channel",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID: catalog.CapabilityCreateVTEXSegmentID,
				Input: capability.Input{
					"brand":         brand,
					"regionId":      regionID,
					"salesChannel":  salesChannel,
					"includeCookie": includeCookie,
				},
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
				Confirmed:     confirmed,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&brand, "brand", "exito", "VTEX brand account to query: exito or carulla")
	command.Flags().StringVar(&regionID, "region-id", "", "VTEX region ID to place in the segment")
	command.Flags().StringVar(&salesChannel, "sales-channel", "", "VTEX sales channel/trade policy to place in the segment")
	command.Flags().BoolVar(&includeCookie, "include-cookie", false, "Include the vtex_segment cookie string in successful output")
	command.Flags().BoolVar(&confirmed, "confirm", false, "Explicitly confirm VTEX segment creation")
	_ = command.MarkFlagRequired("region-id")
	_ = command.MarkFlagRequired("sales-channel")
	return command
}

func newCatalogIntelligentSearchCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "intelligent-search",
		Short: "Run VTEX Intelligent Search commands",
	}
	command.AddCommand(newCatalogIntelligentSearchProductsCommand(bootstrap, options))
	command.AddCommand(newCatalogIntelligentSearchRegionalizedProductsCommand(bootstrap, options))
	return command
}

func newCatalogIntelligentSearchProductsCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var brand string
	var tradePolicy string
	var text string
	var by string
	var values []string
	var query string
	var facets []string
	var page int
	var count int
	var sort string
	var locale string
	var hideUnavailable bool
	var includeUnavailable bool
	var simulationBehavior string
	var cookies []string

	command := &cobra.Command{
		Use:   "products",
		Short: "Search VTEX Intelligent Search products",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			input := capability.Input{
				"brand":       brand,
				"tradePolicy": tradePolicy,
				"page":        page,
				"count":       count,
			}
			if text != "" {
				input["text"] = text
			}
			if by != "" {
				input["by"] = by
			}
			if len(values) > 0 {
				input["value"] = values
			}
			if query != "" {
				input["query"] = query
			}
			if len(facets) > 0 {
				input["facet"] = facets
			}
			if sort != "" {
				input["sort"] = sort
			}
			if locale != "" {
				input["locale"] = locale
			}
			if hideUnavailable || includeUnavailable {
				input["hideUnavailable"] = hideUnavailable && !includeUnavailable
			}
			if simulationBehavior != "" {
				input["simulationBehavior"] = simulationBehavior
			}
			if len(cookies) > 0 {
				input["cookie"] = cookies
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID:  catalog.CapabilityIntelligentSearchProductsID,
				Input:         input,
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&brand, "brand", "exito", "VTEX brand account to query: exito or carulla")
	command.Flags().StringVar(&tradePolicy, "trade-policy", "", "Required VTEX trade policy/sales channel")
	command.Flags().StringVar(&text, "text", "", "Natural-language search text")
	command.Flags().StringVar(&by, "by", "", "Typed lookup mode: product-id, sku-id, ean, sku-reference, slug, or id")
	command.Flags().StringArrayVar(&values, "value", nil, "Lookup value used with --by; repeat for same-type multi-ID lookup")
	command.Flags().StringVar(&query, "query", "", "Raw Intelligent Search query expression")
	command.Flags().StringArrayVar(&facets, "facet", nil, "Additional path facet as key=value; repeat for multiple facets")
	command.Flags().IntVar(&page, "page", 1, "Result page")
	command.Flags().IntVar(&count, "count", 24, "Products per page")
	command.Flags().StringVar(&sort, "sort", "", "Sort value, such as price:asc or orders:desc; empty means relevance")
	command.Flags().StringVar(&locale, "locale", "", "BCP 47 locale")
	command.Flags().BoolVar(&hideUnavailable, "hide-unavailable", false, "Hide unavailable products")
	command.Flags().BoolVar(&includeUnavailable, "include-unavailable", false, "Include unavailable products")
	command.Flags().StringVar(&simulationBehavior, "simulation-behavior", "", "Simulation behavior: default, skip, or only1P")
	command.Flags().StringArrayVar(&cookies, "cookie", nil, "VTEX cookie string for advanced diagnostics; values are redacted from output")
	_ = command.MarkFlagRequired("trade-policy")
	command.MarkFlagsMutuallyExclusive("hide-unavailable", "include-unavailable")
	command.MarkFlagsMutuallyExclusive("text", "query")
	command.MarkFlagsMutuallyExclusive("text", "by")
	command.MarkFlagsMutuallyExclusive("query", "by")
	return command
}

func newCatalogIntelligentSearchRegionalizedProductsCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var brand string
	var country string
	var tradePolicy string
	var longitude string
	var latitude string
	var text string
	var by string
	var values []string
	var query string
	var facets []string
	var page int
	var count int
	var sort string
	var locale string
	var hideUnavailable bool
	var includeUnavailable bool
	var simulationBehavior string
	var confirmed bool

	command := &cobra.Command{
		Use:   "regionalized-products",
		Short: "Resolve a VTEX region, create a segment, and search Intelligent Search products",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			input := capability.Input{
				"brand":       brand,
				"country":     country,
				"tradePolicy": tradePolicy,
				"longitude":   longitude,
				"latitude":    latitude,
				"page":        page,
				"count":       count,
			}
			if text != "" {
				input["text"] = text
			}
			if by != "" {
				input["by"] = by
			}
			if len(values) > 0 {
				input["value"] = values
			}
			if query != "" {
				input["query"] = query
			}
			if len(facets) > 0 {
				input["facet"] = facets
			}
			if sort != "" {
				input["sort"] = sort
			}
			if locale != "" {
				input["locale"] = locale
			}
			if hideUnavailable || includeUnavailable {
				input["hideUnavailable"] = hideUnavailable && !includeUnavailable
			}
			if simulationBehavior != "" {
				input["simulationBehavior"] = simulationBehavior
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID:  workflow.CapabilityRegionalizedIntelligentSearchProductsID,
				Input:         input,
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
				Confirmed:     confirmed,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&brand, "brand", "exito", "VTEX brand account to query: exito or carulla")
	command.Flags().StringVar(&country, "country", "COL", "Country code for VTEX Checkout Regions")
	command.Flags().StringVar(&tradePolicy, "trade-policy", "", "Required VTEX trade policy/sales channel")
	command.Flags().StringVar(&longitude, "longitude", "", "Longitude used in VTEX geoCoordinates")
	command.Flags().StringVar(&latitude, "latitude", "", "Latitude used in VTEX geoCoordinates")
	command.Flags().StringVar(&text, "text", "", "Natural-language search text")
	command.Flags().StringVar(&by, "by", "", "Typed lookup mode: product-id, sku-id, ean, sku-reference, slug, or id")
	command.Flags().StringArrayVar(&values, "value", nil, "Lookup value used with --by; repeat for same-type multi-ID lookup")
	command.Flags().StringVar(&query, "query", "", "Raw Intelligent Search query expression")
	command.Flags().StringArrayVar(&facets, "facet", nil, "Additional path facet as key=value; repeat for multiple facets")
	command.Flags().IntVar(&page, "page", 1, "Result page")
	command.Flags().IntVar(&count, "count", 24, "Products per page")
	command.Flags().StringVar(&sort, "sort", "", "Sort value, such as price:asc or orders:desc; empty means relevance")
	command.Flags().StringVar(&locale, "locale", "", "BCP 47 locale")
	command.Flags().BoolVar(&hideUnavailable, "hide-unavailable", false, "Hide unavailable products")
	command.Flags().BoolVar(&includeUnavailable, "include-unavailable", false, "Include unavailable products")
	command.Flags().StringVar(&simulationBehavior, "simulation-behavior", "", "Simulation behavior: default, skip, or only1P")
	command.Flags().BoolVar(&confirmed, "confirm", false, "Explicitly confirm VTEX segment creation for regionalized search")
	_ = command.MarkFlagRequired("trade-policy")
	_ = command.MarkFlagRequired("longitude")
	_ = command.MarkFlagRequired("latitude")
	command.MarkFlagsMutuallyExclusive("hide-unavailable", "include-unavailable")
	command.MarkFlagsMutuallyExclusive("text", "query")
	command.MarkFlagsMutuallyExclusive("text", "by")
	command.MarkFlagsMutuallyExclusive("query", "by")
	return command
}

func newGeoCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "geo",
		Short: "Run Geo Domain commands",
	}
	command.AddCommand(newGeoGeocodeAddressCommand(bootstrap, options))
	command.AddCommand(newGeoResolveVTEXRegionCommand(bootstrap, options))
	return command
}

func newGeoGeocodeAddressCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var city string
	var address string

	command := &cobra.Command{
		Use:   "geocode-address",
		Short: "Geocode a city/address pair",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID: geo.CapabilityGeocodeAddressID,
				Input: capability.Input{
					"city":    city,
					"address": address,
				},
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&city, "city", "", "City name accepted by the Geo provider")
	command.Flags().StringVar(&address, "address", "", "Address to geocode")
	_ = command.MarkFlagRequired("city")
	_ = command.MarkFlagRequired("address")
	return command
}

func newGeoResolveVTEXRegionCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var brand string
	var country string
	var salesChannel string
	var longitude string
	var latitude string

	command := &cobra.Command{
		Use:   "resolve-vtex-region",
		Short: "Resolve VTEX region coverage from known coordinates",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID: geo.CapabilityResolveVTEXRegionID,
				Input: capability.Input{
					"brand":        brand,
					"country":      country,
					"salesChannel": salesChannel,
					"longitude":    longitude,
					"latitude":     latitude,
				},
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&brand, "brand", "exito", "VTEX brand account to query: exito or carulla")
	command.Flags().StringVar(&country, "country", "COL", "VTEX country code")
	command.Flags().StringVar(&salesChannel, "sales-channel", "", "VTEX sales channel/trade policy passed as sc")
	command.Flags().StringVar(&longitude, "longitude", "", "Longitude coordinate; sent before latitude in geoCoordinates")
	command.Flags().StringVar(&latitude, "latitude", "", "Latitude coordinate; sent after longitude in geoCoordinates")
	_ = command.MarkFlagRequired("sales-channel")
	_ = command.MarkFlagRequired("longitude")
	_ = command.MarkFlagRequired("latitude")
	return command
}

func newOrdersCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "orders",
		Short: "Run Orders Domain commands",
	}
	command.AddCommand(newOrdersGetCommand(bootstrap, options))
	command.AddCommand(newOrdersGetVTEXCommand(bootstrap, options))
	return command
}

func newOrdersGetCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var orderID string
	var orderType string

	command := &cobra.Command{
		Use:   "get",
		Short: "Get a GEOMS order by ID",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID: orders.CapabilityGetID,
				Input: capability.Input{
					"id":        orderID,
					"orderType": orderType,
				},
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&orderID, "id", "", "Order identifier")
	command.Flags().StringVar(&orderType, "order-type", "ExitoEcomm", "GEOMS order type filter, such as ExitoEcomm or CarullaEcomm")
	_ = command.MarkFlagRequired("id")
	return command
}

func newOrdersGetVTEXCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var orderID string
	var brand string

	command := &cobra.Command{
		Use:   "get-vtex",
		Short: "Get a VTEX OMS order by ID",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID: orders.CapabilityGetVTEXID,
				Input: capability.Input{
					"id":    orderID,
					"brand": brand,
				},
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&orderID, "id", "", "VTEX OMS order identifier")
	command.Flags().StringVar(&brand, "brand", "exito", "VTEX brand account to query: exito or carulla")
	_ = command.MarkFlagRequired("id")
	return command
}

func newRunCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var inputJSON string
	var inputFile string
	var confirmed bool

	command := &cobra.Command{
		Use:   "run <capability-id>",
		Short: "Run a capability by its stable ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := parseRunInput(cmd, inputJSON, inputFile)
			if err != nil {
				return err
			}

			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID:  args[0],
				Input:         input,
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
				Confirmed:     confirmed,
			})
			if err != nil {
				return err
			}

			return writeExecutionEnvelope(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&inputJSON, "input-json", "", "Complete capability input object as inline JSON")
	command.Flags().StringVar(&inputFile, "input-file", "", "Path to a JSON file containing the complete capability input object")
	command.Flags().BoolVar(&confirmed, "confirm", false, "Explicitly confirm a confirmation-required capability")
	return command
}

func writeExecutionEnvelope(w io.Writer, envelope capability.Envelope[any]) error {
	if err := presenter.WriteJSON(w, envelope); err != nil {
		return err
	}
	if !envelope.OK {
		return ExitError{Code: ExitCodeFailure}
	}
	return nil
}

func parseRunInput(cmd *cobra.Command, inputJSON string, inputFile string) (capability.Input, error) {
	sources := 0
	if inputJSON != "" {
		sources++
	}
	if inputFile != "" {
		sources++
	}
	stdinAvailable := runStdinAvailable(cmd)
	if stdinAvailable {
		sources++
	}
	if sources > 1 {
		return nil, fmt.Errorf("run input must be provided by only one source")
	}

	switch {
	case inputJSON != "":
		return decodeRunInput([]byte(inputJSON))
	case inputFile != "":
		content, err := os.ReadFile(inputFile) // #nosec G304 -- users explicitly choose the generic run input file
		if err != nil {
			return nil, err
		}
		return decodeRunInput(content)
	case stdinAvailable:
		content, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return nil, err
		}
		return decodeRunInput(content)
	default:
		return capability.Input{}, nil
	}
}

func decodeRunInput(content []byte) (capability.Input, error) {
	var input capability.Input
	if err := json.Unmarshal(content, &input); err != nil {
		return nil, err
	}
	if input == nil {
		return nil, fmt.Errorf("run input must be a JSON object")
	}
	return input, nil
}

func runStdinAvailable(cmd *cobra.Command) bool {
	input := cmd.InOrStdin()
	if input == nil {
		return false
	}
	if input != os.Stdin {
		return true
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return stat.Mode()&os.ModeCharDevice == 0
}

func appOptions(options rootOptions) app.Options {
	return app.Options{
		Config: config.Options{
			ConfigPath: options.configPath,
			Profile:    options.profile,
		},
	}
}

func rootLong(registeredEntries int) string {
	return fmt.Sprintf(
		"Exito Tools command-line interface\n\nExito Tools is the machine-first CLI surface for the application.\n\nRegistered capabilities: %d\n\nUse an implemented subcommand for machine-readable JSON output.",
		registeredEntries,
	)
}
