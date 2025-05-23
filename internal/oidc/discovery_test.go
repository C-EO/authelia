package oidc_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
	"github.com/authelia/authelia/v4/internal/oidc"
)

func TestNewOpenIDConnectProviderDiscovery(t *testing.T) {
	provider := oidc.NewOpenIDConnectProvider(&schema.Configuration{
		IdentityProviders: schema.IdentityProviders{
			OIDC: &schema.IdentityProvidersOpenIDConnect{
				IssuerCertificateChain:   schema.X509CertificateChain{},
				IssuerPrivateKey:         x509PrivateKeyRSA2048,
				HMACSecret:               "asbdhaaskmdlkamdklasmdlkams",
				EnablePKCEPlainChallenge: true,
				Clients: []schema.IdentityProvidersOpenIDConnectClient{
					{
						ID:                  "a-client",
						Secret:              tOpenIDConnectPlainTextClientSecret,
						AuthorizationPolicy: onefactor,
						RedirectURIs: []string{
							"https://google.com",
						},
					},
				},
			},
		},
	}, nil, nil)

	a := provider.GetOpenIDConnectWellKnownConfiguration("https://auth.example.com")

	data, err := json.Marshal(&a)
	assert.NoError(t, err)

	b := oidc.OpenIDConnectWellKnownConfiguration{}

	assert.NoError(t, json.Unmarshal(data, &b))

	assert.Equal(t, a, b)

	y := provider.GetOAuth2WellKnownConfiguration("https://auth.example.com")

	data, err = json.Marshal(&y)
	assert.NoError(t, err)

	z := oidc.OAuth2WellKnownConfiguration{}

	assert.NoError(t, json.Unmarshal(data, &z))

	assert.Equal(t, y, z)
}

func TestNewOpenIDConnectProvider_GetOpenIDConnectWellKnownConfiguration(t *testing.T) {
	provider := oidc.NewOpenIDConnectProvider(&schema.Configuration{
		IdentityProviders: schema.IdentityProviders{
			OIDC: &schema.IdentityProvidersOpenIDConnect{
				IssuerCertificateChain: schema.X509CertificateChain{},
				IssuerPrivateKey:       x509PrivateKeyRSA2048,
				HMACSecret:             "asbdhaaskmdlkamdklasmdlkams",
				Clients: []schema.IdentityProvidersOpenIDConnectClient{
					{
						ID:                  "a-client",
						Secret:              tOpenIDConnectPlainTextClientSecret,
						AuthorizationPolicy: onefactor,
						RedirectURIs: []string{
							"https://google.com",
						},
					},
				},
			},
		},
	}, nil, nil)

	require.NotNil(t, provider)

	disco := provider.GetOpenIDConnectWellKnownConfiguration(examplecom)

	assert.Equal(t, examplecom, disco.Issuer)
	assert.Equal(t, "https://example.com/jwks.json", disco.JWKSURI)
	assert.Equal(t, "https://example.com/api/oidc/authorization", disco.AuthorizationEndpoint)
	assert.Equal(t, "https://example.com/api/oidc/token", disco.TokenEndpoint)
	assert.Equal(t, "https://example.com/api/oidc/userinfo", disco.UserinfoEndpoint)
	assert.Equal(t, "https://example.com/api/oidc/introspection", disco.IntrospectionEndpoint)
	assert.Equal(t, "https://example.com/api/oidc/revocation", disco.RevocationEndpoint)
	assert.Equal(t, "", disco.RegistrationEndpoint)

	assert.Len(t, disco.CodeChallengeMethodsSupported, 1)
	assert.Contains(t, disco.CodeChallengeMethodsSupported, oidc.PKCEChallengeMethodSHA256)

	assert.Len(t, disco.ScopesSupported, 7)
	assert.Contains(t, disco.ScopesSupported, oidc.ScopeOpenID)
	assert.Contains(t, disco.ScopesSupported, oidc.ScopeOfflineAccess)
	assert.Contains(t, disco.ScopesSupported, oidc.ScopeProfile)
	assert.Contains(t, disco.ScopesSupported, oidc.ScopeGroups)
	assert.Contains(t, disco.ScopesSupported, oidc.ScopeEmail)
	assert.Contains(t, disco.ScopesSupported, oidc.ScopeAddress)
	assert.Contains(t, disco.ScopesSupported, oidc.ScopePhone)

	assert.Len(t, disco.ResponseModesSupported, 7)
	assert.Contains(t, disco.ResponseModesSupported, oidc.ResponseModeFormPost)
	assert.Contains(t, disco.ResponseModesSupported, oidc.ResponseModeQuery)
	assert.Contains(t, disco.ResponseModesSupported, oidc.ResponseModeFragment)
	assert.Contains(t, disco.ResponseModesSupported, oidc.ResponseModeJWT)
	assert.Contains(t, disco.ResponseModesSupported, oidc.ResponseModeFormPostJWT)
	assert.Contains(t, disco.ResponseModesSupported, oidc.ResponseModeQueryJWT)
	assert.Contains(t, disco.ResponseModesSupported, oidc.ResponseModeFragmentJWT)

	assert.Len(t, disco.SubjectTypesSupported, 2)
	assert.Contains(t, disco.SubjectTypesSupported, oidc.SubjectTypePublic)
	assert.Contains(t, disco.SubjectTypesSupported, oidc.SubjectTypePairwise)

	assert.Len(t, disco.ResponseTypesSupported, 7)
	assert.Contains(t, disco.ResponseTypesSupported, oidc.ResponseTypeAuthorizationCodeFlow)
	assert.Contains(t, disco.ResponseTypesSupported, oidc.ResponseTypeImplicitFlowIDToken)
	assert.Contains(t, disco.ResponseTypesSupported, oidc.ResponseTypeImplicitFlowToken)
	assert.Contains(t, disco.ResponseTypesSupported, oidc.ResponseTypeImplicitFlowBoth)
	assert.Contains(t, disco.ResponseTypesSupported, oidc.ResponseTypeHybridFlowIDToken)
	assert.Contains(t, disco.ResponseTypesSupported, oidc.ResponseTypeHybridFlowToken)
	assert.Contains(t, disco.ResponseTypesSupported, oidc.ResponseTypeHybridFlowBoth)

	assert.Len(t, disco.TokenEndpointAuthMethodsSupported, 5)
	assert.Contains(t, disco.TokenEndpointAuthMethodsSupported, oidc.ClientAuthMethodClientSecretBasic)
	assert.Contains(t, disco.TokenEndpointAuthMethodsSupported, oidc.ClientAuthMethodClientSecretPost)
	assert.Contains(t, disco.TokenEndpointAuthMethodsSupported, oidc.ClientAuthMethodClientSecretJWT)
	assert.Contains(t, disco.TokenEndpointAuthMethodsSupported, oidc.ClientAuthMethodPrivateKeyJWT)
	assert.Contains(t, disco.TokenEndpointAuthMethodsSupported, oidc.ClientAuthMethodNone)

	assert.Len(t, disco.RevocationEndpointAuthMethodsSupported, 5)
	assert.Contains(t, disco.RevocationEndpointAuthMethodsSupported, oidc.ClientAuthMethodClientSecretBasic)
	assert.Contains(t, disco.RevocationEndpointAuthMethodsSupported, oidc.ClientAuthMethodClientSecretPost)
	assert.Contains(t, disco.RevocationEndpointAuthMethodsSupported, oidc.ClientAuthMethodClientSecretJWT)
	assert.Contains(t, disco.RevocationEndpointAuthMethodsSupported, oidc.ClientAuthMethodPrivateKeyJWT)
	assert.Contains(t, disco.RevocationEndpointAuthMethodsSupported, oidc.ClientAuthMethodNone)

	assert.Equal(t, []string{oidc.ClientAuthMethodClientSecretBasic, oidc.ClientAuthMethodClientSecretPost, oidc.ClientAuthMethodClientSecretJWT, oidc.ClientAuthMethodPrivateKeyJWT}, disco.IntrospectionEndpointAuthMethodsSupported)
	assert.Equal(t, []string{oidc.GrantTypeAuthorizationCode, oidc.GrantTypeImplicit, oidc.GrantTypeClientCredentials, oidc.GrantTypeRefreshToken, oidc.GrantTypeDeviceCode}, disco.GrantTypesSupported)
	assert.Equal(t, []string{oidc.SigningAlgHMACUsingSHA256, oidc.SigningAlgHMACUsingSHA384, oidc.SigningAlgHMACUsingSHA512, oidc.SigningAlgRSAUsingSHA256, oidc.SigningAlgRSAUsingSHA384, oidc.SigningAlgRSAUsingSHA512, oidc.SigningAlgECDSAUsingP256AndSHA256, oidc.SigningAlgECDSAUsingP384AndSHA384, oidc.SigningAlgECDSAUsingP521AndSHA512, oidc.SigningAlgRSAPSSUsingSHA256, oidc.SigningAlgRSAPSSUsingSHA384, oidc.SigningAlgRSAPSSUsingSHA512}, disco.RevocationEndpointAuthSigningAlgValuesSupported)
	assert.Equal(t, []string{oidc.SigningAlgHMACUsingSHA256, oidc.SigningAlgHMACUsingSHA384, oidc.SigningAlgHMACUsingSHA512, oidc.SigningAlgRSAUsingSHA256, oidc.SigningAlgRSAUsingSHA384, oidc.SigningAlgRSAUsingSHA512, oidc.SigningAlgECDSAUsingP256AndSHA256, oidc.SigningAlgECDSAUsingP384AndSHA384, oidc.SigningAlgECDSAUsingP521AndSHA512, oidc.SigningAlgRSAPSSUsingSHA256, oidc.SigningAlgRSAPSSUsingSHA384, oidc.SigningAlgRSAPSSUsingSHA512}, disco.TokenEndpointAuthSigningAlgValuesSupported)
	assert.Equal(t, []string{oidc.SigningAlgHMACUsingSHA256, oidc.SigningAlgHMACUsingSHA384, oidc.SigningAlgHMACUsingSHA512, oidc.SigningAlgRSAUsingSHA256, oidc.SigningAlgRSAUsingSHA384, oidc.SigningAlgRSAUsingSHA512, oidc.SigningAlgECDSAUsingP256AndSHA256, oidc.SigningAlgECDSAUsingP384AndSHA384, oidc.SigningAlgECDSAUsingP521AndSHA512, oidc.SigningAlgRSAPSSUsingSHA256, oidc.SigningAlgRSAPSSUsingSHA384, oidc.SigningAlgRSAPSSUsingSHA512}, disco.IDTokenSigningAlgValuesSupported)
	assert.Equal(t, []string{oidc.SigningAlgHMACUsingSHA256, oidc.SigningAlgHMACUsingSHA384, oidc.SigningAlgHMACUsingSHA512, oidc.SigningAlgRSAUsingSHA256, oidc.SigningAlgRSAUsingSHA384, oidc.SigningAlgRSAUsingSHA512, oidc.SigningAlgECDSAUsingP256AndSHA256, oidc.SigningAlgECDSAUsingP384AndSHA384, oidc.SigningAlgECDSAUsingP521AndSHA512, oidc.SigningAlgRSAPSSUsingSHA256, oidc.SigningAlgRSAPSSUsingSHA384, oidc.SigningAlgRSAPSSUsingSHA512, oidc.SigningAlgNone}, disco.UserinfoSigningAlgValuesSupported)
	assert.Equal(t, []string{oidc.SigningAlgHMACUsingSHA256, oidc.SigningAlgHMACUsingSHA384, oidc.SigningAlgHMACUsingSHA512, oidc.SigningAlgRSAUsingSHA256, oidc.SigningAlgRSAUsingSHA384, oidc.SigningAlgRSAUsingSHA512, oidc.SigningAlgECDSAUsingP256AndSHA256, oidc.SigningAlgECDSAUsingP384AndSHA384, oidc.SigningAlgECDSAUsingP521AndSHA512, oidc.SigningAlgRSAPSSUsingSHA256, oidc.SigningAlgRSAPSSUsingSHA384, oidc.SigningAlgRSAPSSUsingSHA512, oidc.SigningAlgNone}, disco.RequestObjectSigningAlgValuesSupported)

	assert.Len(t, disco.ClaimsSupported, 33)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimAuthenticationMethodsReference)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimAudience)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimAuthorizedParty)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimClientIdentifier)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimExpirationTime)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimIssuedAt)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimIssuer)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimJWTID)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimRequestedAt)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimSubject)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimAuthenticationTime)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimNonce)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimEmailAlts)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimGroups)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimFullName)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimGivenName)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimFamilyName)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimMiddleName)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimNickname)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimPreferredUsername)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimProfile)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimPicture)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimWebsite)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimEmail)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimEmailVerified)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimGender)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimBirthdate)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimZoneinfo)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimLocale)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimPhoneNumber)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimPhoneNumberVerified)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimAddress)
	assert.Contains(t, disco.ClaimsSupported, oidc.ClaimUpdatedAt)

	assert.Len(t, disco.PromptValuesSupported, 4)
	assert.Contains(t, disco.PromptValuesSupported, oidc.PromptConsent)
	assert.Contains(t, disco.PromptValuesSupported, oidc.PromptSelectAccount)
	assert.Contains(t, disco.PromptValuesSupported, oidc.PromptLogin)
	assert.Contains(t, disco.PromptValuesSupported, oidc.PromptNone)
}

func TestNewOpenIDConnectWellKnownConfiguration_Copy(t *testing.T) {
	config := &oidc.OpenIDConnectWellKnownConfiguration{
		OAuth2WellKnownConfiguration: oidc.OAuth2WellKnownConfiguration{
			CommonDiscoveryOptions: oidc.CommonDiscoveryOptions{
				Issuer:                            "https://example.com",
				JWKSURI:                           "https://example.com/jwks.json",
				AuthorizationEndpoint:             "",
				TokenEndpoint:                     "",
				SubjectTypesSupported:             nil,
				ResponseTypesSupported:            nil,
				GrantTypesSupported:               nil,
				ResponseModesSupported:            nil,
				ScopesSupported:                   nil,
				ClaimsSupported:                   nil,
				UILocalesSupported:                nil,
				TokenEndpointAuthMethodsSupported: nil,
				TokenEndpointAuthSigningAlgValuesSupported: nil,
				ServiceDocumentation:                       "",
				OPPolicyURI:                                "",
				OPTOSURI:                                   "",
			},
			OAuth2DiscoveryOptions: oidc.OAuth2DiscoveryOptions{
				IntrospectionEndpoint:                              "",
				RevocationEndpoint:                                 "",
				RegistrationEndpoint:                               "",
				IntrospectionEndpointAuthMethodsSupported:          nil,
				RevocationEndpointAuthMethodsSupported:             nil,
				RevocationEndpointAuthSigningAlgValuesSupported:    nil,
				IntrospectionEndpointAuthSigningAlgValuesSupported: nil,
				CodeChallengeMethodsSupported:                      nil,
			},
			OAuth2DeviceAuthorizationGrantDiscoveryOptions: &oidc.OAuth2DeviceAuthorizationGrantDiscoveryOptions{
				DeviceAuthorizationEndpoint: "",
			},
			OAuth2MutualTLSClientAuthenticationDiscoveryOptions: &oidc.OAuth2MutualTLSClientAuthenticationDiscoveryOptions{
				TLSClientCertificateBoundAccessTokens: false,
				MutualTLSEndpointAliases: oidc.OAuth2MutualTLSClientAuthenticationAliasesDiscoveryOptions{
					AuthorizationEndpoint:              "",
					TokenEndpoint:                      "",
					IntrospectionEndpoint:              "",
					RevocationEndpoint:                 "",
					EndSessionEndpoint:                 "",
					UserinfoEndpoint:                   "",
					BackChannelAuthenticationEndpoint:  "",
					FederationRegistrationEndpoint:     "",
					PushedAuthorizationRequestEndpoint: "",
					RegistrationEndpoint:               "",
				},
			},
			OAuth2IssuerIdentificationDiscoveryOptions: &oidc.OAuth2IssuerIdentificationDiscoveryOptions{
				AuthorizationResponseIssuerParameterSupported: false,
			},
			OAuth2JWTIntrospectionResponseDiscoveryOptions: &oidc.OAuth2JWTIntrospectionResponseDiscoveryOptions{
				IntrospectionSigningAlgValuesSupported:    nil,
				IntrospectionEncryptionAlgValuesSupported: nil,
				IntrospectionEncryptionEncValuesSupported: nil,
			},
			OAuth2JWTSecuredAuthorizationRequestDiscoveryOptions: &oidc.OAuth2JWTSecuredAuthorizationRequestDiscoveryOptions{
				RequireSignedRequestObject: false,
			},
			OAuth2PushedAuthorizationDiscoveryOptions: &oidc.OAuth2PushedAuthorizationDiscoveryOptions{
				PushedAuthorizationRequestEndpoint: "",
				RequirePushedAuthorizationRequests: false,
			},
		},
		OpenIDConnectDiscoveryOptions: oidc.OpenIDConnectDiscoveryOptions{
			UserinfoEndpoint:                          "",
			IDTokenSigningAlgValuesSupported:          nil,
			UserinfoSigningAlgValuesSupported:         nil,
			RequestObjectSigningAlgValuesSupported:    nil,
			IDTokenEncryptionAlgValuesSupported:       nil,
			UserinfoEncryptionAlgValuesSupported:      nil,
			RequestObjectEncryptionAlgValuesSupported: nil,
			IDTokenEncryptionEncValuesSupported:       nil,
			UserinfoEncryptionEncValuesSupported:      nil,
			RequestObjectEncryptionEncValuesSupported: nil,
			ACRValuesSupported:                        nil,
			DisplayValuesSupported:                    nil,
			ClaimTypesSupported:                       nil,
			ClaimLocalesSupported:                     nil,
			RequestParameterSupported:                 true,
			RequestURIParameterSupported:              true,
			RequireRequestURIRegistration:             true,
			ClaimsParameterSupported:                  true,
		},
		OpenIDConnectFrontChannelLogoutDiscoveryOptions: &oidc.OpenIDConnectFrontChannelLogoutDiscoveryOptions{
			FrontChannelLogoutSupported:        false,
			FrontChannelLogoutSessionSupported: false,
		},
		OpenIDConnectBackChannelLogoutDiscoveryOptions: &oidc.OpenIDConnectBackChannelLogoutDiscoveryOptions{
			BackChannelLogoutSupported:        false,
			BackChannelLogoutSessionSupported: false,
		},
		OpenIDConnectSessionManagementDiscoveryOptions: &oidc.OpenIDConnectSessionManagementDiscoveryOptions{
			CheckSessionIFrame: "",
		},
		OpenIDConnectRPInitiatedLogoutDiscoveryOptions: &oidc.OpenIDConnectRPInitiatedLogoutDiscoveryOptions{
			EndSessionEndpoint: "",
		},
		OpenIDConnectPromptCreateDiscoveryOptions: &oidc.OpenIDConnectPromptCreateDiscoveryOptions{
			PromptValuesSupported: nil,
		},
		OpenIDConnectClientInitiatedBackChannelAuthFlowDiscoveryOptions: &oidc.OpenIDConnectClientInitiatedBackChannelAuthFlowDiscoveryOptions{
			BackChannelAuthenticationEndpoint:               "",
			BackChannelTokenDeliveryModesSupported:          nil,
			BackChannelAuthRequestSigningAlgValuesSupported: nil,
			BackChannelUserCodeParameterSupported:           false,
		},
		OpenIDConnectJWTSecuredAuthorizationResponseModeDiscoveryOptions: &oidc.OpenIDConnectJWTSecuredAuthorizationResponseModeDiscoveryOptions{
			AuthorizationSigningAlgValuesSupported:    nil,
			AuthorizationEncryptionAlgValuesSupported: nil,
			AuthorizationEncryptionEncValuesSupported: nil,
		},
		OpenIDFederationDiscoveryOptions: &oidc.OpenIDFederationDiscoveryOptions{
			FederationRegistrationEndpoint:                 "",
			ClientRegistrationTypesSupported:               nil,
			RequestAuthenticationMethodsSupported:          nil,
			RequestAuthenticationSigningAlgValuesSupported: nil,
		},
	}

	x := config.Copy()

	assert.Equal(t, config, &x)

	y := config.OAuth2WellKnownConfiguration.Copy()

	assert.Equal(t, config.OAuth2WellKnownConfiguration, y)
}
