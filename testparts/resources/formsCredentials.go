package resources

var GoogleFormCredentials = `
	{
		"installed": {
			"client_id": "579206520663-gq882tt82a9v7ctsu2dhibjulpvp6srf.apps.googleusercontent.com",
			"project_id": "classtracker",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_secret": "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",
			"redirect_uris": [ "http://localhost:%d/auth/google/callback" ]
		}
	}
`
