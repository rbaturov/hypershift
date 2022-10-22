package powervs

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/hypershift/cmd/cluster/core"
	powervsinfra "github.com/openshift/hypershift/cmd/infra/powervs"
	"github.com/openshift/hypershift/cmd/log"
)

func NewDestroyCommand(opts *core.DestroyOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "powervs",
		Short:        "Destroys a HostedCluster and its resources on PowerVS",
		SilenceUsage: true,
	}

	opts.PowerVSPlatform = core.PowerVSPlatformDestroyOptions{
		Region:    "us-south",
		Zone:      "us-south",
		VPCRegion: "us-south",
	}

	cmd.Flags().StringVar(&opts.PowerVSPlatform.ResourceGroup, "resource-group", opts.PowerVSPlatform.ResourceGroup, "IBM Cloud Resource group")
	cmd.Flags().StringVar(&opts.InfraID, "infra-id", opts.InfraID, "Cluster ID with which to tag IBM Cloud resources")
	cmd.Flags().StringVar(&opts.PowerVSPlatform.BaseDomain, "base-domain", opts.PowerVSPlatform.BaseDomain, "Cluster's base domain")
	cmd.Flags().StringVar(&opts.PowerVSPlatform.Region, "region", opts.PowerVSPlatform.Region, "IBM Cloud region. Default is us-south")
	cmd.Flags().StringVar(&opts.PowerVSPlatform.Zone, "zone", opts.PowerVSPlatform.Zone, "IBM Cloud zone. Default is us-south")
	cmd.Flags().StringVar(&opts.PowerVSPlatform.VPCRegion, "vpc-region", opts.PowerVSPlatform.VPCRegion, "IBM Cloud VPC Region for VPC resources. Default is us-south")
	cmd.Flags().BoolVar(&opts.PowerVSPlatform.Debug, "debug", opts.PowerVSPlatform.Debug, "Enabling this will print PowerVS API Request & Response logs")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)
		go func() {
			<-sigs
			cancel()
		}()

		if err := DestroyCluster(ctx, opts); err != nil {
			log.Log.Error(err, "Failed to destroy cluster")
			os.Exit(1)
		}
	}

	return cmd
}

func DestroyCluster(ctx context.Context, o *core.DestroyOptions) error {
	hostedCluster, err := core.GetCluster(ctx, o)
	if err != nil {
		return err
	}
	if hostedCluster != nil {
		o.InfraID = hostedCluster.Spec.InfraID
		o.PowerVSPlatform.BaseDomain = hostedCluster.Spec.DNS.BaseDomain
		o.PowerVSPlatform.ResourceGroup = hostedCluster.Spec.Platform.PowerVS.ResourceGroup
		o.PowerVSPlatform.Region = hostedCluster.Spec.Platform.PowerVS.Region
		o.PowerVSPlatform.Zone = hostedCluster.Spec.Platform.PowerVS.Zone
		o.PowerVSPlatform.VPCRegion = hostedCluster.Spec.Platform.PowerVS.VPC.Region
		o.PowerVSPlatform.CISCRN = hostedCluster.Spec.Platform.PowerVS.CISInstanceCRN
		o.PowerVSPlatform.CISDomainID = hostedCluster.Spec.DNS.PrivateZoneID
	}

	var inputErrors []error
	if o.InfraID == "" {
		inputErrors = append(inputErrors, fmt.Errorf("infrastructure ID is required"))
	}
	if o.PowerVSPlatform.BaseDomain == "" {
		inputErrors = append(inputErrors, fmt.Errorf("base domain is required"))
	}
	if o.PowerVSPlatform.Region == "" {
		inputErrors = append(inputErrors, fmt.Errorf("PowerVS region is required"))
	}
	if o.PowerVSPlatform.Zone == "" {
		inputErrors = append(inputErrors, fmt.Errorf("PowerVS zone is required"))
	}
	if o.PowerVSPlatform.VPCRegion == "" {
		inputErrors = append(inputErrors, fmt.Errorf("VPC region is required"))
	}
	if o.PowerVSPlatform.ResourceGroup == "" {
		inputErrors = append(inputErrors, fmt.Errorf("resource group is required"))
	}
	if err := errors.NewAggregate(inputErrors); err != nil {
		return fmt.Errorf("required inputs are missing: %w", err)
	}

	return core.DestroyCluster(ctx, hostedCluster, o, destroyPlatformSpecifics)
}

func destroyPlatformSpecifics(ctx context.Context, o *core.DestroyOptions) error {
	return (&powervsinfra.DestroyInfraOptions{
		Name:          o.Name,
		Namespace:     o.Namespace,
		InfraID:       o.InfraID,
		BaseDomain:    o.PowerVSPlatform.BaseDomain,
		CISCRN:        o.PowerVSPlatform.CISCRN,
		CISDomainID:   o.PowerVSPlatform.CISDomainID,
		ResourceGroup: o.PowerVSPlatform.ResourceGroup,
		Region:        o.PowerVSPlatform.Region,
		Zone:          o.PowerVSPlatform.Zone,
		VPCRegion:     o.PowerVSPlatform.VPCRegion,
		Debug:         o.PowerVSPlatform.Debug,
	}).Run(ctx)
}
