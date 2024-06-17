package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/charmbracelet/log"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"os"
	"os/exec"
)

const ACTLABS_APP_ID string = "bee16ca1-a401-40ee-bb6a-34349ebd993e"
const GROUP_NAME string = "repro-project"

type AzureValues struct {
	TenantID       string
	SubscriptionID string
	Location       string
	User           string
	UserID         string
}

type AzureCLIAccountOutput struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}

func ProvisionLabsInfrastructure(tenant string, sub string, location string) {
	ctx := context.Background()

	// get our values squared away before we go too far
	azureValues := AzureValues{
		TenantID:       tenant,
		SubscriptionID: sub,
		Location:       location,
	}

	if azureValues.TenantID == "" || azureValues.SubscriptionID == "" {
		log.Debug("No tenant and/or subscription ID provided - attempting to source them from Azure CLI.")
		pullValuesFromCLI(ctx, &azureValues)
	}

	// Provision the labs infrastructure
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal("Can't get credentials for performing Azure operations - exiting. Error: ", err.Error())
	}

	// grab the user's UPN
	graphClient, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(cred, []string{"Files.Read"})
	if err != nil {
		log.Fatal("Error creating Graph client: ", err.Error())
	}

	res, err := graphClient.Me().Get(ctx, nil)
	if err != nil {
		log.Fatal("Error getting user from Graph: ", err.Error())
	}
	azureValues.User = *res.GetUserPrincipalName()

	log.Info("Checking access levels on the subscription to ensure the right permissions are in place...")
	// Check if the current user has the Owner role on the subscription
	// If not, we'll need to prompt the user to add the role
	authzFactory, err := armauthorization.NewClientFactory(azureValues.SubscriptionID, cred, nil)
	if err != nil {
		log.Fatal("Error creating authorization client factory: ", err.Error(), ". We can't proceed without knowing permissions.")
		os.Exit(1)
	}

	// parse the list of role assignments on the subscription. filter down to the Owner role and see if our current user is in there
	pager := authzFactory.NewRoleAssignmentsClient().NewListForSubscriptionPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatal("Error getting role assignments: ", err.Error())
			os.Exit(1)
		}
		res := page.RoleAssignmentListResult.Value
		for _, assignment := range res {

			if *assignment.Properties.PrincipalID == azureValues.UserID &&
				*assignment.Properties.RoleDefinitionID == fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635", azureValues.SubscriptionID) {
				log.Info("User has Owner permissions on subscription - proceeding with provisioning.")
				break
			}
		}
	}

	// At this point, we've got all the right permissions, so away we go!
	// Create a resource group
	log.Info("Creating resource group...")
	rgClient, err := armresources.NewResourceGroupsClient(azureValues.SubscriptionID, cred, nil)
	if err != nil {
		log.Fatal("Error creating resource group client: ", err.Error())
		os.Exit(1)
	}

	params := armresources.ResourceGroup{
		Location: &azureValues.Location,
	}
	rg, err = rgClient.CreateOrUpdate(ctx, GROUP_NAME, params, nil)
	if err != nil {
		log.Fatal("Error creating resource group: ", err.Error())
		os.Exit(1)
	}
	
}

func pullValuesFromCLI(ctx context.Context, a *AzureValues) {
	// Execute `az account show` and retrieve the tenant and subscription IDs from the resulting JSON
	cmd := exec.Command("az", "account", "show")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Error running 'az account show': ", err.Error(), " - unable to proceed with provisioning.")
		os.Exit(1)
	}

	var accountOutput AzureCLIAccountOutput
	err = json.Unmarshal(output, &accountOutput)
	if err != nil {
		log.Fatal("Error unmarshalling 'az account show' output: ", err.Error(), " - unable to proceed with provisioning.")
		os.Exit(1)
	}

	a.TenantID = accountOutput.TenantID
	a.SubscriptionID = accountOutput.ID

}
