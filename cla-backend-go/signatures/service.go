// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/communitybridge/easycla/cla-backend-go/github"
	"github.com/communitybridge/easycla/cla-backend-go/github_organizations"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/users"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v1/restapi/operations/signatures"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v1/models"
	githubpkg "github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

// SignatureService interface
type SignatureService interface {
	GetSignature(ctx context.Context, signatureID string) (*models.Signature, error)
	GetIndividualSignature(ctx context.Context, claGroupID, userID string, approved, signed *bool) (*models.Signature, error)
	GetCorporateSignature(ctx context.Context, claGroupID, companyID string, approved, signed *bool) (*models.Signature, error)
	GetProjectSignatures(ctx context.Context, params signatures.GetProjectSignaturesParams) (*models.Signatures, error)
	CreateProjectSummaryReport(ctx context.Context, params signatures.CreateProjectSummaryReportParams) (*models.SignatureReport, error)
	GetProjectCompanySignature(ctx context.Context, companyID, projectID string, approved, signed *bool, nextKey *string, pageSize *int64) (*models.Signature, error)
	GetProjectCompanySignatures(ctx context.Context, params signatures.GetProjectCompanySignaturesParams) (*models.Signatures, error)
	GetProjectCompanyEmployeeSignatures(ctx context.Context, params signatures.GetProjectCompanyEmployeeSignaturesParams, criteria *ApprovalCriteria) (*models.Signatures, error)
	GetCompanySignatures(ctx context.Context, params signatures.GetCompanySignaturesParams) (*models.Signatures, error)
	GetCompanyIDsWithSignedCorporateSignatures(ctx context.Context, claGroupID string) ([]SignatureCompanyID, error)
	GetUserSignatures(ctx context.Context, params signatures.GetUserSignaturesParams) (*models.Signatures, error)
	InvalidateProjectRecords(ctx context.Context, projectID, note string) (int, error)

	GetGithubOrganizationsFromApprovalList(ctx context.Context, signatureID string, githubAccessToken string) ([]models.GithubOrg, error)
	AddGithubOrganizationToApprovalList(ctx context.Context, signatureID string, approvalListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error)
	DeleteGithubOrganizationFromApprovalList(ctx context.Context, signatureID string, approvalListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error)
	UpdateApprovalList(ctx context.Context, authUser *auth.User, claGroupModel *models.ClaGroup, companyModel *models.Company, claGroupID string, params *models.ApprovalList) (*models.Signature, error)

	AddCLAManager(ctx context.Context, signatureID, claManagerID string) (*models.Signature, error)
	RemoveCLAManager(ctx context.Context, ignatureID, claManagerID string) (*models.Signature, error)

	GetClaGroupICLASignatures(ctx context.Context, claGroupID string, searchTerm *string, approved, signed *bool, pageSize int64, nextKey string) (*models.IclaSignatures, error)
	GetClaGroupCCLASignatures(ctx context.Context, claGroupID string, approved, signed *bool) (*models.Signatures, error)
	GetClaGroupCorporateContributors(ctx context.Context, claGroupID string, companyID *string, searchTerm *string) (*models.CorporateContributorList, error)
}

type service struct {
	repo                SignatureRepository
	companyService      company.IService
	usersService        users.Service
	eventsService       events.Service
	githubOrgValidation bool
	repositoryService   repositories.Service
	githubOrgService    github_organizations.ServiceInterface
}

// NewService creates a new signature service
func NewService(repo SignatureRepository, companyService company.IService, usersService users.Service, eventsService events.Service, githubOrgValidation bool, repositoryService repositories.Service, githubOrgService github_organizations.ServiceInterface) SignatureService {
	return service{
		repo,
		companyService,
		usersService,
		eventsService,
		githubOrgValidation,
		repositoryService,
		githubOrgService,
	}
}

// GetSignature returns the signature associated with the specified signature ID
func (s service) GetSignature(ctx context.Context, signatureID string) (*models.Signature, error) {
	return s.repo.GetSignature(ctx, signatureID)
}

// GetIndividualSignature returns the signature associated with the specified CLA Group and User ID
func (s service) GetIndividualSignature(ctx context.Context, claGroupID, userID string, approved, signed *bool) (*models.Signature, error) {
	return s.repo.GetIndividualSignature(ctx, claGroupID, userID, approved, signed)
}

// GetCorporateSignature returns the signature associated with the specified CLA Group and Company ID
func (s service) GetCorporateSignature(ctx context.Context, claGroupID, companyID string, approved, signed *bool) (*models.Signature, error) {
	return s.repo.GetCorporateSignature(ctx, claGroupID, companyID, approved, signed)
}

// GetProjectSignatures returns the list of signatures associated with the specified project
func (s service) GetProjectSignatures(ctx context.Context, params signatures.GetProjectSignaturesParams) (*models.Signatures, error) {

	projectSignatures, err := s.repo.GetProjectSignatures(ctx, params)
	if err != nil {
		return nil, err
	}

	return projectSignatures, nil
}

// CreateProjectSummaryReport generates a project summary report based on the specified input
func (s service) CreateProjectSummaryReport(ctx context.Context, params signatures.CreateProjectSummaryReportParams) (*models.SignatureReport, error) {

	projectSignatures, err := s.repo.CreateProjectSummaryReport(ctx, params)
	if err != nil {
		return nil, err
	}

	return projectSignatures, nil
}

// GetProjectCompanySignature returns the signature associated with the specified project and company
func (s service) GetProjectCompanySignature(ctx context.Context, companyID, projectID string, approved, signed *bool, nextKey *string, pageSize *int64) (*models.Signature, error) {
	return s.repo.GetProjectCompanySignature(ctx, companyID, projectID, approved, signed, nextKey, pageSize)
}

// GetProjectCompanySignatures returns the list of signatures associated with the specified project
func (s service) GetProjectCompanySignatures(ctx context.Context, params signatures.GetProjectCompanySignaturesParams) (*models.Signatures, error) {

	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	signed := true
	approved := true

	projectSignatures, err := s.repo.GetProjectCompanySignatures(
		ctx, params.CompanyID, params.ProjectID, &signed, &approved, params.NextKey, params.SortOrder, &pageSize)
	if err != nil {
		return nil, err
	}

	return projectSignatures, nil
}

// GetProjectCompanyEmployeeSignatures returns the list of employee signatures associated with the specified project
func (s service) GetProjectCompanyEmployeeSignatures(ctx context.Context, params signatures.GetProjectCompanyEmployeeSignaturesParams, criteria *ApprovalCriteria) (*models.Signatures, error) {

	if params.PageSize == nil {
		params.PageSize = utils.Int64(10)
	}

	projectSignatures, err := s.repo.GetProjectCompanyEmployeeSignatures(ctx, params, criteria)
	if err != nil {
		return nil, err
	}

	return projectSignatures, nil
}

// GetCompanySignatures returns the list of signatures associated with the specified company
func (s service) GetCompanySignatures(ctx context.Context, params signatures.GetCompanySignaturesParams) (*models.Signatures, error) {

	const defaultPageSize int64 = 50
	var pageSize = defaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	companySignatures, err := s.repo.GetCompanySignatures(ctx, params, pageSize, LoadACLDetails)
	if err != nil {
		return nil, err
	}

	return companySignatures, nil
}

// GetCompanyIDsWithSignedCorporateSignatures returns a list of company IDs that have signed a CLA agreement
func (s service) GetCompanyIDsWithSignedCorporateSignatures(ctx context.Context, claGroupID string) ([]SignatureCompanyID, error) {
	return s.repo.GetCompanyIDsWithSignedCorporateSignatures(ctx, claGroupID)
}

// GetUserSignatures returns the list of user signatures associated with the specified user
func (s service) GetUserSignatures(ctx context.Context, params signatures.GetUserSignaturesParams) (*models.Signatures, error) {

	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	userSignatures, err := s.repo.GetUserSignatures(ctx, params, pageSize)
	if err != nil {
		return nil, err
	}

	return userSignatures, nil
}

// GetGithubOrganizationsFromApprovalList retrieves the organization from the approval list
func (s service) GetGithubOrganizationsFromApprovalList(ctx context.Context, signatureID string, githubAccessToken string) ([]models.GithubOrg, error) {

	if signatureID == "" {
		msg := "unable to get GitHub organizations approval list - signature ID is nil"
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	orgIds, err := s.repo.GetGithubOrganizationsFromApprovalList(ctx, signatureID)
	if err != nil {
		log.Warnf("error loading github organization from approval list using signatureID: %s, error: %v",
			signatureID, err)
		return nil, err
	}

	if githubAccessToken != "" {
		log.Debugf("already authenticated with github - scanning for user's orgs...")

		selectedOrgs := make(map[string]struct{}, len(orgIds))
		for _, selectedOrg := range orgIds {
			selectedOrgs[*selectedOrg.ID] = struct{}{}
		}

		// Since we're logged into github, lets get the list of organization we can add.
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubAccessToken},
		)
		tc := oauth2.NewClient(utils.NewContext(), ts)
		client := githubpkg.NewClient(tc)

		opt := &githubpkg.ListOptions{
			PerPage: 100,
		}

		orgs, _, err := client.Organizations.List(utils.NewContext(), "", opt)
		if err != nil {
			return nil, err
		}

		for _, org := range orgs {
			_, ok := selectedOrgs[*org.Login]
			if ok {
				continue
			}

			orgIds = append(orgIds, models.GithubOrg{ID: org.Login})
		}
	}

	return orgIds, nil
}

// AddGithubOrganizationToApprovalList adds the GH organization to the approval list
func (s service) AddGithubOrganizationToApprovalList(ctx context.Context, signatureID string, approvalListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error) {
	organizationID := approvalListParams.OrganizationID

	if signatureID == "" {
		msg := "unable to add GitHub organization from approval list - signature ID is nil"
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	if organizationID == nil {
		msg := "unable to add GitHub organization from approval list - organization ID is nil"
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	// GH_ORG_VALIDATION environment - set to false to test locally which will by-pass the GH auth checks and
	// allow functional tests (e.g. with curl or postmon) - default is enabled

	if s.githubOrgValidation {
		// Verify the authenticated github user has access to the github organization being added.
		if githubAccessToken == "" {
			msg := fmt.Sprintf("unable to add github organization, not logged in using "+
				"signatureID: %s, github organization id: %s",
				signatureID, *organizationID)
			log.Warn(msg)
			return nil, errors.New(msg)
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubAccessToken},
		)
		tc := oauth2.NewClient(utils.NewContext(), ts)
		client := githubpkg.NewClient(tc)

		opt := &githubpkg.ListOptions{
			PerPage: 100,
		}

		log.Debugf("querying for user's github organizations...")
		orgs, _, err := client.Organizations.List(utils.NewContext(), "", opt)
		if err != nil {
			return nil, err
		}

		found := false
		for _, org := range orgs {
			if *org.Login == *organizationID {
				found = true
				break
			}
		}

		if !found {
			msg := fmt.Sprintf("user is not authorized for github organization id: %s", *organizationID)
			log.Warnf(msg)
			return nil, errors.New(msg)
		}
	}

	gitHubOrgApprovalList, err := s.repo.AddGithubOrganizationToApprovalList(ctx, signatureID, *organizationID)
	if err != nil {
		log.Warnf("issue adding github organization to approval list using signatureID: %s, gh org id: %s, error: %v",
			signatureID, *organizationID, err)
		return nil, err
	}

	return gitHubOrgApprovalList, nil
}

// DeleteGithubOrganizationFromApprovalList deletes the specified GH organization from the approval list
func (s service) DeleteGithubOrganizationFromApprovalList(ctx context.Context, signatureID string, approvalListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error) {

	// Extract the payload values
	organizationID := approvalListParams.OrganizationID

	if signatureID == "" {
		msg := "unable to delete GitHub organization from approval list - signature ID is nil"
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	if organizationID == nil {
		msg := "unable to delete GitHub organization from approval list - organization ID is nil"
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	// GH_ORG_VALIDATION environment - set to false to test locally which will by-pass the GH auth checks and
	// allow functional tests (e.g. with curl or postmon) - default is enabled

	if s.githubOrgValidation {
		// Verify the authenticated github user has access to the github organization being added.
		if githubAccessToken == "" {
			msg := fmt.Sprintf("unable to delete github organization, not logged in using "+
				"signatureID: %s, github organization id: %s",
				signatureID, *organizationID)
			log.Warn(msg)
			return nil, errors.New(msg)
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubAccessToken},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client := githubpkg.NewClient(tc)

		opt := &githubpkg.ListOptions{
			PerPage: 100,
		}

		log.Debugf("querying for user's github organizations...")
		orgs, _, err := client.Organizations.List(context.Background(), "", opt)
		if err != nil {
			return nil, err
		}

		found := false
		for _, org := range orgs {
			if *org.Login == *organizationID {
				found = true
				break
			}
		}

		if !found {
			msg := fmt.Sprintf("user is not authorized for github organization id: %s", *organizationID)
			log.Warnf(msg)
			return nil, errors.New(msg)
		}
	}

	gitHubOrgApprovalList, err := s.repo.DeleteGithubOrganizationFromApprovalList(ctx, signatureID, *organizationID)
	if err != nil {
		return nil, err
	}

	return gitHubOrgApprovalList, nil
}

// UpdateApprovalList service method which handles updating the various approval lists
func (s service) UpdateApprovalList(ctx context.Context, authUser *auth.User, claGroupModel *models.ClaGroup, companyModel *models.Company, claGroupID string, params *models.ApprovalList) (*models.Signature, error) { // nolint gocyclo
	f := logrus.Fields{
		"functionName":      "v1.signatures.service.UpdateApprovalList",
		utils.XREQUESTID:    ctx.Value(utils.XREQUESTID),
		"authUser.UserName": authUser.UserName,
		"authUser.Email":    authUser.Email,
		"claGroupID":        claGroupID,
		"claGroupName":      claGroupModel.ProjectName,
		"companyName":       companyModel.CompanyName,
		"companyID":         companyModel.CompanyID,
	}

	log.WithFields(f).Debugf("processing update approval list request")

	// Lookup the project corporate signature - should have one
	pageSize := int64(1)
	signed, approved := true, true
	corporateSigModel, sigErr := s.GetProjectCompanySignature(ctx, companyModel.CompanyID, claGroupID, &signed, &approved, nil, &pageSize)
	if sigErr != nil {
		msg := fmt.Sprintf("unable to locate project company signature by Company ID: %s, Project ID: %s, CLA Group ID: %s, error: %+v",
			companyModel.CompanyID, claGroupModel.ProjectID, claGroupID, sigErr)
		log.WithFields(f).WithError(sigErr).Warn(msg)
		return nil, NewBadRequestError(msg)
	}
	// If not found, return error
	if corporateSigModel == nil {
		msg := fmt.Sprintf("unable to locate signature for company ID: %s CLA Group ID: %s, type: ccla, signed: %t, approved: %t",
			companyModel.CompanyID, claGroupID, signed, approved)
		log.WithFields(f).Warn(msg)
		return nil, NewBadRequestError(msg)
	}

	// Ensure current user is in the Signature ACL
	claManagers := corporateSigModel.SignatureACL
	if !utils.CurrentUserInACL(authUser, claManagers) {
		msg := fmt.Sprintf("EasyCLA - 403 Forbidden - CLA Manager %s / %s is not authorized to approve request for company ID: %s / %s / %s, project ID: %s / %s / %s",
			authUser.UserName, authUser.Email,
			companyModel.CompanyName, companyModel.CompanyExternalID, companyModel.CompanyID,
			claGroupModel.ProjectName, claGroupModel.ProjectExternalID, claGroupModel.ProjectID)
		return nil, NewForbiddenError(msg)
	}

	// Lookup the user making the request - should be the CLA Manager
	userModel, userErr := s.usersService.GetUserByUserName(authUser.UserName, true)
	if userErr != nil {
		log.WithFields(f).WithError(userErr).Warnf("unable to lookup user by user name: %s", authUser.UserName)
		return nil, userErr
	}

	// This event is ONLY used when we need to invalidate the signature
	eventArgs := &events.LogEventArgs{
		EventType:     events.InvalidatedSignature, // reviewed and
		ProjectID:     claGroupModel.ProjectExternalID,
		ClaGroupModel: claGroupModel,
		CompanyID:     companyModel.CompanyID,
		CompanyModel:  companyModel,
		LfUsername:    userModel.LfUsername,
		UserID:        userModel.UserID,
		UserModel:     userModel,
		ProjectSFID:   claGroupModel.ProjectExternalID,
	}

	// Here we perform the approval list updates for all the different types of approval lists
	log.WithFields(f).Debugf("updating approval list...")
	updatedSig, err := s.repo.UpdateApprovalList(ctx, userModel, claGroupModel, companyModel.CompanyID, params, eventArgs)
	if err != nil {
		log.WithFields(f).WithError(err).Warnf("problem updating approval list for company ID: %s, project ID: %s, cla group ID: %s", companyModel.CompanyID, claGroupModel.ProjectID, claGroupID)
		return updatedSig, err
	}

	// Log Events that the CLA manager updated the approval lists
	log.WithFields(f).Debugf("creating event log entry...")
	go s.createEventLogEntries(ctx, companyModel, claGroupModel, userModel, params)

	// Send an email to each of the CLA Managers
	log.WithFields(f).Debugf("sending email to cla managers...")
	for _, claManager := range claManagers {
		claManagerEmail := getBestEmail(&claManager) // nolint
		s.sendApprovalListUpdateEmailToCLAManagers(companyModel, claGroupModel, claManager.Username, claManagerEmail, params)
	}

	// TODO: DAD - update email template to indicate that if auto crate ECLA is enabled, that users should be good-to-go
	// Send emails to contributors if email or GitHub/GitLab username was added or removed
	log.WithFields(f).Debugf("sending email to contributors...")
	s.sendRequestAccessEmailToContributors(authUser, companyModel, claGroupModel, params)

	// If auto create ECLA is enabled for this Corporate Agreement, then create an ECLA for each employee that was added to the approval list
	// TODO: DAD should we move this to above the actual approval list update and email blast?
	log.WithFields(f).Debugf("checking for auto-create ECLA option: %t...", corporateSigModel.AutoCreateECLA)
	if corporateSigModel.AutoCreateECLA {
		createdECLARecord := false
		log.WithFields(f).Debugf("auto-create ECLA option is enabled: %t...", corporateSigModel.AutoCreateECLA)

		// For the add email list, create an ECLA signature record for each user
		var employeeUserModel *models.User
		var userLookupErr error
		for _, email := range params.AddEmailApprovalList {
			log.WithFields(f).Debugf("auto-create ECLA option - add email: %s", email)

			// Lookup the user by email in the local EasyCLA database - this will exist if the user first
			// initiated the request from GitHub and if they shared their email (made it public). This record will
			// likely not exist if the CLA Manager added the email directly from the UI without the user first
			// initiating the workflow.
			employeeUserModel, userLookupErr = s.usersService.GetUserByEmail(email)
			// If we couldn't find the user, then create a user record
			if userLookupErr != nil || employeeUserModel == nil {
				log.WithFields(f).WithError(userLookupErr).Warnf("unable to lookup existing user by email: %s", email)
				var userCreateErr error
				// Create a new user record based on the email and company ID
				employeeUserModel, userCreateErr = s.createUserModel("", "", "", "", email, companyModel.CompanyID, "auto-create ECLA user from CLA Manager approval list update")
				if userCreateErr != nil || employeeUserModel == nil {
					log.WithFields(f).WithError(userCreateErr).Warnf("unable to create a new user with email: %s", email)
					// TODO: DAD - how do we communicate this back to the CLA Manager in the UI - simply return the error?
					return nil, userCreateErr
				}
			}

			// Ok, auto-create the employee acknowledgement record
			createErr := s.repo.CreateProjectCompanyEmployeeSignature(ctx, companyModel, claGroupModel, employeeUserModel)
			if createErr != nil {
				log.WithFields(f).WithError(createErr).Warnf("unable to create project company employee signature record for: %+v", employeeUserModel)
				// TODO: DAD - how do we communicate this back to the CLA Manager in the UI - simply return the error?
				return nil, createErr
			}
			createdECLARecord = true
		}
		for _, gitHubUserName := range params.AddGithubUsernameApprovalList {
			log.WithFields(f).Debugf("auto-create ECLA option - add githubUserName: %s", gitHubUserName)

			// Lookup the user by GitHub username in the local EasyCLA database - this will exist if the user first
			// initiated the request from GitHub. This record will likely not exist if the CLA Manager added the GitHub
			// username directly from the UI without the user first initiating the workflow.
			employeeUserModel, userLookupErr = s.usersService.GetUserByGitHubUsername(gitHubUserName)
			// If we couldn't find the user, then create a user record
			if userLookupErr != nil || employeeUserModel == nil {
				log.WithFields(f).WithError(userLookupErr).Infof("unable to lookup existing user by GitHub username: %s in our local database - will attempt to create a new record", gitHubUserName)
				var gitHubUserID = ""
				var gitHubUserEmail = ""
				// Attempt to lookup the GitHub user record by the GitHub username - we need the GitHub numeric ID value which was not provided by the UI/API call
				gitHubUserModel, gitHubErr := github.GetUserDetails(gitHubUserName)
				// Should get a model, no errors and have at least the ID
				if gitHubErr != nil || gitHubUserModel == nil || gitHubUserModel.ID == nil {
					log.WithFields(f).WithError(gitHubErr).Warnf("problem looking up GitHub user details for user: %s, model: %+v, error: %+v", gitHubUserName, gitHubUserModel, gitHubErr)
					// TODO: DAD - should we fail the entire operation? - seems like we should abort if we can't lookup the user details - maybe we should add GitHub and GitLab validation to the UI or early in the request?
					return nil, gitHubErr
				}

				if gitHubUserModel.ID != nil {
					gitHubUserID = strconv.FormatInt(*gitHubUserModel.ID, 10)
				}
				// User may not have a public email
				if gitHubUserModel.Email != nil {
					gitHubUserEmail = *gitHubUserModel.Email
				}

				var userCreateErr error
				// Create a new user record based on the GitHub information, email and company ID
				employeeUserModel, userCreateErr = s.createUserModel(gitHubUserName, gitHubUserID, "", "", gitHubUserEmail, companyModel.CompanyID, "auto-create ECLA user from CLA Manager approval list update")
				if userCreateErr != nil || employeeUserModel == nil {
					log.WithFields(f).WithError(userCreateErr).Warnf("unable to create a new user with GitHub username: %s", gitHubUserName)
					// TODO: DAD - how do we communicate this back to the CLA Manager in the UI - simply return the error?
					return nil, userCreateErr
				}
			}

			// Ok, finally, auto-create the employee acknowledgement record
			createErr := s.repo.CreateProjectCompanyEmployeeSignature(ctx, companyModel, claGroupModel, employeeUserModel)
			if createErr != nil {
				log.WithFields(f).WithError(createErr).Warnf("unable to create project company employee signature record for: %+v", employeeUserModel)
				// TODO: DAD - how do we communicate this back to the CLA Manager in the UI - simply return the error?
				return nil, createErr
			}
			createdECLARecord = true
		}

		/* Note: GitLab API is currently not working - plus, what credentials do we use to lookup the user details? Our current API leverages the GitHub app credentials (again, broken as of 09/2022 due to needing to refresh the token every hour)
		for _, gitLabUserName := range params.AddGitlabUsernameApprovalList {
			// Lookup the user by GitLab username in the local EasyCLA database - this will exist if the user first
			// initiated the request from GitLab. This record will likely not exist if the CLA Manager added the GitLab
			// username directly from the UI without the user first initiating the workflow.
			employeeUserModel, userLookupErr := s.usersService.GetUserByGitLabUsername(gitLabUserName)
			// If we couldn't find the user, then create a user record
			if userLookupErr != nil || employeeUserModel == nil {
				log.WithFields(f).WithError(userLookupErr).Warnf("unable to lookup existing user by GitLab username: %s", gitLabUserName)
				var gitLabUserID = ""
				var gitLabUserEmail = ""
				// GitLab API is currently not working - plus, what credentials do we use to lookup the user details? Our current API leverages the GitHub app credentials (again, broken as of 09/2022 due to needing to refresh the token every hour)
				gitHubUserModel, gitHubErr := gitlab.List(gitLabUserName)
				// Should get a model, no errors and have at least the ID
				if gitHubErr != nil || gitHubUserModel == nil || gitHubUserModel.ID == nil {
					log.WithFields(f).WithError(gitHubErr).Warnf("problem looking up GitHub user details for user: %s, model: %+v, error: %+v", gitLabUserName, gitHubUserModel, gitHubErr)
					// TODO: DAD - should we fail the entire operation? - seems like we should abort if we can't lookup the user details - maybe we should add GitHub and GitLab validation to the UI or early in the request?
					return nil, gitHubErr
				} else {
					if gitHubUserModel.ID != nil {
						gitLabUserID = strconv.FormatInt(*gitHubUserModel.ID, 10)
					}
					// User may not have a public email
					if gitHubUserModel.Email != nil {
						gitLabUserEmail = *gitHubUserModel.Email
					}
				}

				var userCreateErr error
				// Create a new user record based on the GitHub information, email and company ID
				employeeUserModel, userCreateErr = s.createUserModel("", "", gitLabUserName, gitLabUserID, gitLabUserEmail, companyModel.CompanyID, "auto-create ECLA user from CLA Manager approval list update")
				if userCreateErr != nil || employeeUserModel == nil {
					log.WithFields(f).WithError(userCreateErr).Warnf("unable to create a new user with GitLab username: %s", gitLabUserName)
					// TODO: DAD - how do we communicate this back to the CLA Manager in the UI - simply return the error?
					return nil, userCreateErr
				}
			}

			// Ok, finally, auto-create the employee acknowledgement record
			createErr := s.repo.CreateProjectCompanyEmployeeSignature(ctx, companyModel, claGroupModel, employeeUserModel)
			if createErr != nil {
				log.WithFields(f).WithError(createErr).Warnf("unable to create project company employee signature record for: %+v", employeeUserModel)
				// TODO: DAD - how do we communicate this back to the CLA Manager in the UI - simply return the error?
				return nil, createErr
			}
			createdECLARecord = true
		}
		*/

		if createdECLARecord && employeeUserModel != nil {
			log.WithFields(f).Debug("created one or more ECLA records - need to update GitHub status check")
			// TODO: add GitHub status check update
			signatureMetadata, sigErr := s.repo.GetActiveSignatureMetadata(ctx, employeeUserModel.UserID)
			if sigErr != nil {
				log.WithFields(f).WithError(sigErr).Warnf("unable to get active signature record for : %+v", employeeUserModel)
				return nil, sigErr
			}

			// Fetch easycla repository
			claRepository, repoErr := s.repositoryService.GetRepository(ctx, signatureMetadata.RepositoryID)
			if repoErr != nil {
				log.WithFields(f).WithError(repoErr).Warnf("unable to fetch repository by ID : %s ", signatureMetadata.RepositoryID)
				return nil, repoErr
			}

			if !claRepository.Enabled {
				log.WithFields(f).Debugf("Repository: %s associated with PR: %s is NOT enabled", claRepository.RepositoryURL, signatureMetadata.PullRequestID)
				return nil, errors.New("Repository is not enabled")
			}

			// fetch gihub org details
			githubOrg, githubOrgErr := s.githubOrgService.GetGitHubOrganizationByName(ctx, claRepository.RepositoryName)
			if githubOrgErr != nil {
				log.WithFields(f).WithError(githubOrgErr).Warn("unable to get githubOrg")
				return nil, githubOrgErr
			}

			repositoryID, idErr := strconv.Atoi(signatureMetadata.RepositoryID)
			if idErr != nil {
				return nil, idErr
			}

			pullRequestID, idErr := strconv.Atoi(signatureMetadata.PullRequestID)
			if idErr != nil {
				return nil, idErr
			}

			// Update change request
			updateErr := s.updateChangeRequest(ctx, githubOrg, int64(repositoryID), int64(pullRequestID), signatureMetadata.ProjectID)
			if updateErr != nil {
				log.WithFields(f).WithError(updateErr).Warnf("unable to update pull request: %d ", pullRequestID)
				return nil, updateErr
			}

		} else {

			log.WithFields(f).Debug("no ECLA records created - no need to update GitHub status check")
		}
	}

	return updatedSig, nil
}

// InvalidateProjectRecords disassociates project signatures
func (s service) InvalidateProjectRecords(ctx context.Context, projectID, note string) (int, error) {
	f := logrus.Fields{
		"functionName": "v1.signatures.service.InvalidateProjectRecords",
		"projectID":    projectID,
	}

	result, err := s.repo.ProjectSignatures(ctx, projectID)
	if err != nil {
		log.WithFields(f).WithError(err).Warnf(fmt.Sprintf("Unable to get signatures for project: %s", projectID))
		return 0, err
	}

	if len(result.Signatures) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(result.Signatures))
		log.WithFields(f).Debugf(fmt.Sprintf("Invalidating %d signatures for project: %s ",
			len(result.Signatures), projectID))
		for _, signature := range result.Signatures {
			// Do this in parallel, as we could have a lot to invalidate
			go func(sigID, projectID string) {
				defer wg.Done()
				updateErr := s.repo.InvalidateProjectRecord(ctx, sigID, note)
				if updateErr != nil {
					log.WithFields(f).WithError(updateErr).Warnf("Unable to update signature: %s with project ID: %s, error: %v", sigID, projectID, updateErr)
				}
			}(signature.SignatureID, projectID)
		}

		// Wait until all the workers are done
		wg.Wait()
	}

	return len(result.Signatures), nil
}

// AddCLAManager adds the specified manager to the signature ACL list
func (s service) AddCLAManager(ctx context.Context, signatureID, claManagerID string) (*models.Signature, error) {
	return s.repo.AddCLAManager(ctx, signatureID, claManagerID)
}

// RemoveCLAManager removes the specified manager from the signature ACL list
func (s service) RemoveCLAManager(ctx context.Context, signatureID, claManagerID string) (*models.Signature, error) {
	return s.repo.RemoveCLAManager(ctx, signatureID, claManagerID)
}

// appendList is a helper function to generate the email content of the Approval List changes
func appendList(approvalList []string, message string) string {
	approvalListSummary := ""

	if len(approvalList) > 0 {
		for _, value := range approvalList {
			approvalListSummary += fmt.Sprintf("<li>%s %s</li>", message, value)
		}
	}

	return approvalListSummary
}

// buildApprovalListSummary is a helper function to generate the email content of the Approval List changes
func buildApprovalListSummary(approvalListChanges *models.ApprovalList) string {
	approvalListSummary := "<ul>"
	approvalListSummary += appendList(approvalListChanges.AddEmailApprovalList, "Added Email:")
	approvalListSummary += appendList(approvalListChanges.RemoveEmailApprovalList, "Removed Email:")
	approvalListSummary += appendList(approvalListChanges.AddDomainApprovalList, "Added Domain:")
	approvalListSummary += appendList(approvalListChanges.RemoveDomainApprovalList, "Removed Domain:")
	approvalListSummary += appendList(approvalListChanges.AddGithubUsernameApprovalList, "Added GitHub User:")
	approvalListSummary += appendList(approvalListChanges.RemoveGithubUsernameApprovalList, "Removed GitHub User:")
	approvalListSummary += appendList(approvalListChanges.AddGithubOrgApprovalList, "Added GitHub Organization:")
	approvalListSummary += appendList(approvalListChanges.RemoveGithubOrgApprovalList, "Removed GitHub Organization:")
	approvalListSummary += appendList(approvalListChanges.AddGitlabUsernameApprovalList, "Added Gitlab User:")
	approvalListSummary += appendList(approvalListChanges.RemoveGitlabUsernameApprovalList, "Removed Gitlab User:")
	approvalListSummary += appendList(approvalListChanges.AddGitlabOrgApprovalList, "Added Gitlab Organization:")
	approvalListSummary += appendList(approvalListChanges.RemoveGitlabOrgApprovalList, "Removed Gitlab Organization:")
	approvalListSummary += "</ul>"
	return approvalListSummary
}

func (s service) GetClaGroupICLASignatures(ctx context.Context, claGroupID string, searchTerm *string, approved, signed *bool, pageSize int64, nextKey string) (*models.IclaSignatures, error) {
	return s.repo.GetClaGroupICLASignatures(ctx, claGroupID, searchTerm, approved, signed, pageSize, nextKey)
}

func (s service) GetClaGroupCCLASignatures(ctx context.Context, claGroupID string, approved, signed *bool) (*models.Signatures, error) {
	pageSize := utils.Int64(1000)
	return s.repo.GetProjectSignatures(ctx, signatures.GetProjectSignaturesParams{
		ClaType:   aws.String(utils.ClaTypeCCLA),
		ProjectID: claGroupID,
		PageSize:  pageSize,
		Approved:  approved,
		Signed:    signed,
	})
}

func (s service) GetClaGroupCorporateContributors(ctx context.Context, claGroupID string, companyID *string, searchTerm *string) (*models.CorporateContributorList, error) {
	return s.repo.GetClaGroupCorporateContributors(ctx, claGroupID, companyID, searchTerm)
}

// updateChangeRequest is a helper function that updates PR upong auto ecla update
func (s service) updateChangeRequest(ctx context.Context, ghOrg *models.GithubOrganization, repositoryID, pullRequestID int64, projectID string) error {
	f := logrus.Fields{
		"functionName":  "v1.signatures.service.updateChangeRequest",
		"repositoryID":  repositoryID,
		"pullRequestID": pullRequestID,
		"projectID":     projectID,
	}

	githubRepository, ghErr := github.GetGitHubRepository(ctx, ghOrg.OrganizationInstallationID, repositoryID)
	if ghErr != nil {
		log.WithFields(f).WithError(ghErr).Warn("unable to get github repository")
		return ghErr
	}

	// Fetch committers
	log.WithFields(f).Debugf("fetching commit authors for PR: %d", pullRequestID)

	authors, _, authorsErr := github.GetPullRequestCommitAuthors(ctx, ghOrg.OrganizationInstallationID, int(pullRequestID), *githubRepository.Owner.Name, *githubRepository.Name)
	if authorsErr != nil {
		log.WithFields(f).WithError(authorsErr).Warnf("unable to get commit authors for PR: %d", pullRequestID)
		return authorsErr
	}

	signed := make([]*github.UserCommitSummary, 0)
	unsigned := make([]*github.UserCommitSummary, 0)

	// triage signed and unsigned users
	for _, userSummary := range authors {
		if !userSummary.IsValid() {
			unsigned = append(unsigned, userSummary)
		}
		user, userErr := s.usersService.GetUserByGitHubUsername(*userSummary.CommitAuthor.Name)
		if userErr != nil {
			unsigned = append(unsigned, userSummary)
			break
		}
		userSigned, signedErr := s.hasUserSigned(ctx, user, projectID)
		if signedErr != nil {
			break
		}
		if userSigned != nil && *userSigned {
			signed = append(signed, userSummary)
		}
	}

	log.WithFields(f).Debugf("User status signed: %+v and missing: %+v", signed, unsigned)

	// update pull request
	// github.UpdatePullRequest(ctx,)

	return nil
}

func (s service) hasUserSigned(ctx context.Context, user *models.User, projectID string) (*bool, error) {
	f := logrus.Fields{
		"functionName": "v1.signatures.service.updateChangeRequest",
		"projectID":    projectID,
		"user":         user,
	}
	var hasSigned bool
	log.WithFields(f).Debugf("checking to see if user has signed an ICLA ")

	approved := true
	signed := true

	// check for ICLA
	signature, sigErr := s.GetIndividualSignature(ctx, projectID, user.UserID, &approved, &signed)
	if sigErr != nil {
		return nil, sigErr
	}

	if signature != nil {
		hasSigned = true
		log.WithFields(f).Debugf("ICLA signature check passed for user: %+v on project : %s", user, projectID)
	} else {
		log.WithFields(f).Debugf("ICLA signature check failed for user: %+v on project : %s", user, projectID)
	}

	// Check for CCLA
	companyID := user.CompanyID

	if companyID != "" {
		// Get employee signature
		ecla, eclaErr := s.GetProjectCompanyEmployeeSignatures(ctx, signatures.GetProjectCompanyEmployeeSignaturesParams{
			CompanyID: companyID,
			ProjectID: projectID,
		}, &ApprovalCriteria{})

		if eclaErr != nil {
			log.WithFields(f).Debugf("Unable to fetch ecla record for company: %s and project: %s", companyID, projectID)
			return nil, eclaErr
		}
		employeeSignature := ecla.Signatures[0]
		// employeeSignature, empErr := s.repo.GetProjectCompanyEmployeeSignature()

		if employeeSignature != nil {
			log.WithFields(f).Debugf("CCLA Signature check - located employee acknowledgement - signature id: %s", employeeSignature.SignatureID)
			// Get ccla signature of company to access whitelist
			cclaSignature, cclaErr := s.GetCorporateSignature(ctx, projectID, companyID, &approved, &signed)
			if cclaErr != nil {
				return nil, cclaErr
			}

			if cclaSignature != nil {
				approved, approvedErr := s.userIsApproved(ctx, user, cclaSignature)
				if approvedErr != nil {
					return nil, approvedErr
				}
				if approved {
					log.WithFields(f).Debugf("user:%s is in the approval list for signature : %s", user.UserID, signature.SignatureID)
					hasSigned = true
				}
			}
		}

	}

	return &hasSigned, nil
}

func (s service) userIsApproved(ctx context.Context, user *models.User, cclaSignature *models.Signature) (bool, error) {
	emails := append(user.Emails, string(user.LfEmail))

	f := logrus.Fields{
		"functionName": "v1.signatures.service.userIsApproved",
	}

	// check email whitelist
	whitelist := cclaSignature.EmailApprovalList
	if len(whitelist) > 0 {
		for _, email := range emails {
			if s.contains(whitelist, strings.ToLower(strings.TrimSpace(email))) {
				return true, nil
			}
		}
	} else {
		log.WithFields(f).Debugf("no whitelist found for ccla: %s", cclaSignature.SignatureID)
	}

	// check domain whitelist
	domainWhitelist := cclaSignature.DomainApprovalList
	if len(domainWhitelist) > 0 {
		matched, err := s.processPattern(emails, domainWhitelist)
		if err != nil {
			return false, err
		}
		if matched != nil && *matched {
			return true, nil
		}
	}

	// check github whitelist
	if user.GithubUsername != "" {
		githubOrgApprovalList := cclaSignature.GithubOrgApprovalList
		if len(githubOrgApprovalList) > 0 {
			log.WithFields(f).Debugf("determining if github user :%s is associated with ant of the github orgs : %+v", user.GithubUsername, githubOrgApprovalList)
		}

		for _, org := range githubOrgApprovalList {
			membership, err := github.GetMembership(ctx, user.GithubUsername, org)
			if err != nil {
				break
			}
			if membership != nil {
				log.WithFields(f).Debugf("found matching github organization: %s for user: %s", org, user.GithubUsername)
				return true, nil
			} else {
				log.WithFields(f).Debugf("user: %s is not in the organization: %s", user.GithubUsername, org)
			}
		}
	}

	return false, nil
}

func (s service) contains(items []string, val string) bool {
	for _, item := range items {
		if val == item {
			return true
		}
	}
	return false
}

func (s service) processPattern(emails []string, patterns []string) (*bool, error) {
	matched := false

	for _, pattern := range patterns {
		if strings.HasPrefix(pattern, "*.") {
			pattern = strings.Replace(pattern, "*.", ".*", -1)
		} else if strings.HasPrefix(pattern, "*") {
			pattern = strings.Replace(pattern, "*", ".*", -1)
		} else if strings.HasPrefix(pattern, ".") {
			pattern = strings.Replace(pattern, ".", ".*", -1)
		}

		preProcessedPattern := fmt.Sprintf("^.*@%s$", pattern)
		compiled, err := regexp.Compile(preProcessedPattern)
		if err != nil {
			return nil, err
		}

		for _, email := range emails {
			if compiled.MatchString(email) {
				matched = true
				break
			}
		}
	}

	return &matched, nil
}
