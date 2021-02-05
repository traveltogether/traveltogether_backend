package errors

import "github.com/gin-gonic/gin"

var (
	MissingAuthenticationInformation = gin.H{"error": "missing_authentication"}
	InternalError                    = gin.H{"error": "internal_error"}
	InvalidRequest                   = gin.H{"error": "invalid_request"}
	NotFound                         = gin.H{"error": "not_found"}
	Forbidden                        = gin.H{"error": "forbidden"}
	InvalidLoginData                 = gin.H{"error": "invalid_login_data"}
	UserAlreadyExists                = gin.H{"error": "user_already_exists"}
	RequestingOwnJourney             = gin.H{"error": "requesting_own_journey"}
	RequestAlreadyAccepted           = gin.H{"error": "request_already_accepted"}
	RequestAlreadyDeclined           = gin.H{"error": "request_already_declined"}
	JourneyHasBeenCancelled          = gin.H{"error": "journey_has_been_cancelled"}
	JourneyAlreadyTookPlace          = gin.H{"error": "journey_already_took_place"}
	RequestsNotOpen                  = gin.H{"error": "requests_not_open"}
	UserNotRequestedToJoin           = gin.H{"error": "user_not_requested_to_join"}
	UserHasNotBeenAccepted           = gin.H{"error": "user_has_not_been_accepted"}
	UserHasNotBeenDeclined           = gin.H{"error": "user_has_not_been_declined"}
	NotRequestedToJoin               = gin.H{"error": "not_requested_to_join"}
	MailAlreadyInUse                 = gin.H{"error": "mail_already_in_use"}
	InvalidMailAddressOrUsername     = gin.H{"error": "invalid_mail_address_or_username"}
)
