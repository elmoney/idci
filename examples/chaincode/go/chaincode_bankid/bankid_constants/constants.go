package bankid_constants

const (

	BoxReqCreatedName				= "boxRequestCreated"
	BoxReqToVerificationUnitName 	= "boxRequestToVerificationUnit"
	BoxReqApprovedName 				= "boxRequestApproved"
	BoxReqRejectedName				= "boxRequestRejected"

	ErrorFormatSerializeObject			= "[%s]Error at serialize object: %s"
	ErrorFormatSerializeContainer		= "[%s]Error at serialize container requests: %s"
	ErrorFormatDeserializeObject		= "[%s]Error at deserialize object: %s"
	ErrorFormatDeserializeContainer		= "[%s]Error at deserialize container requests : %s"
	ErrorFormatSavingState				= "[%s]Error at putting object to state: %s"
	ErrorFormatGettingState				= "[%s]Error at getting object from state: %s"
	ErrorFormatEncryptObjectInternal 	= "[%s]Error at encrypt object: %s"
	ErrorFormatDecryptObjectInternal 	= "[%s]Error at decrypt object: %s"
	ErrorFormatDecryptObjectCommon 		= "%s_DECRYPT_ERROR"
	ErrorFormatHexDecodeInternal 		= "[%s]Error at %s hex decoding from string: %s"
	ErrorFormatHexDecodeCommon 			= "%s_DECODE_ERROR"

	UnitPubKeysStatePrefix = "pubkey_"
)