// Package classification Trading Simulator API.
//
// the purpose of this parser is to provide an parser
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//	Schemes: http
//	Host: localhost
//	BasePath: /sim
//	Version: 0.1.0
//	License: None
//	Contact: Vyacheslav<WWRadkov@gmail.com> https://gitlab.com/Muhmr
//
//	Consumes:
//	- parser/json
//
//	Produces:
//	- parser/json
//
//	Security:
//	- api_key:
//
//	SecurityDefinitions:
//	api_key:
//	     type: apiKey
//	     name: KEY
//	     in: header
//	oauth2:
//	    type: oauth2
//	    authorizationUrl: /oauth2/auth
//	    tokenUrl: /oauth2/token
//	    in: header
//	    scopes:
//	      bar: foo
//	    flow: accessCode
//
//	Extensions:
//	x-meta-value: value
//	x-meta-array:
//	  - value1
//	  - value2
//	x-meta-array-obj:
//	  - name: obj
//	    value: field
//
// swagger:meta
package docs

import (
	"github.com/emptywe/trading_sim/entity"
	"github.com/emptywe/trading_sim/internal/simulator/api/v1/router"
)

// swagger:route POST /auth/sign-up foobar-tag idOfsignUpEndpoint
// signUp creates user.
// responses:
//   200: signUpResponse
//   422: signUpErrResponse

// If everything is ok creates new user and give back his id.
// swagger:response signUpResponse
type signUpResponseWrapper struct {
	// in:body
	Body entity.SignUpResponse
}

// If user already exist return an error
// swagger:response signUpErrResponse
type signUpNewResponseWrapper struct {
	// in:body
	Body entity.ErrorResponse
}

// swagger:parameters idOfsignUpEndpoint
type signUpParamsWrapper struct {
	// Enter email, username and password to create user.
	// in:body
	Body entity.SignUpRequest
}

// swagger:route POST /auth/sign-in foobar-tag idOfsignInEndpoint
// signIn singing in user.
// responses:
//   200: signInResponse
//   401: signInErrResponse

// If everything is ok signing in new user give back his id and create session_cache with cookie, giving back session_cache id and JWT in response.
// swagger:response signInResponse
type signInResponseWrapper struct {
	// in:body
	Body entity.SignInResponse
}

// If user doesn't exist return an error
// swagger:response signInErrResponse
type signInNewResponseWrapper struct {
	// in:body
	Body entity.ErrorResponse
}

// swagger:parameters idOfsignInEndpoint
type signInParamsWrapper struct {
	// Enter username and password to sign in user.
	// in:body
	Body entity.SignInRequest
}

// swagger:route POST /auth/logout foobar-tag idOflogOutEndpoint
// logOut logouts user.
// responses:
//   200: logOutResponse
//   401: logOutErrResponse

// If everything is ok logging out user and delete session_cache
// swagger:response logOutResponse
type logOutResponseWrapper struct {
	// in:body
	Body struct{}
}

// If user doesn't sign in return status 401
// swagger:response logOutErrResponse
type logOutNewResponseWrapper struct {
	// in:body
	Body struct{}
}

// swagger:parameters idOflogOutEndpoint
type logOutParamsWrapper struct {
	// Enter username and password to sign in user.
	// in:body
	Body struct{}
}

// swagger:route GET /basket/prices foobar-tag idOfpricesEndpoint
// prices show all available currencies.
// responses:
//   200: pricesResponse
//   403: pricesErrResponse

// If everything is ok gives back slice of currencies struct
// swagger:response pricesResponse
type pricesResponseWrapper struct {
	// in:body
	Body [2]entity.CurrencyOutput
}

// If user doesn't sign in return status 403
// swagger:response pricesErrResponse
type pricesNewResponseWrapper struct {
	// in:body
	Body struct{}
}

// swagger:parameters idOfpricesEndpoint
type pricesParamsWrapper struct {
	// Just reach endpoint to get info.
	// in:body
	Body struct{}
}

// swagger:route GET /basket/balance foobar-tag idOfbalanceEndpoint
// balance show users currencies and it's USD amount.
// responses:
//   200: balanceResponse
//   403: balanceErrResponse

// If everything is ok give back slice of currencies struct and sum of all currencies in USD
// swagger:response balanceResponse
type balanceResponseWrapper struct {
	// in:body
	Body entity.BalanceResponse
}

// If user doesn't sign in return status 403
// swagger:response balanceErrResponse
type balanceNewResponseWrapper struct {
	// in:body
	Body struct{}
}

// swagger:parameters idOfbalanceEndpoint
type balanceParamsWrapper struct {
	// Just reach endpoint to get info.
	// in:body
	Body struct{}
}

// swagger:route GET /basket/top foobar-tag idOftopUsersEndpoint
// topUsers show top 10 users
// responses:
//   200: topUsersResponse
//   403: topUsersErrResponse

// If everything is ok give back top 10 users
// swagger:response topUsersResponse
type topUsersResponseWrapper struct {
	// in:body
	Body []entity.TUser
}

// If user doesn't sign in return status 403
// swagger:response topUsersErrResponse
type topUsersNewResponseWrapper struct {
	// in:body
	Body struct{}
}

// swagger:parameters idOftopUsersEndpoint
type topUsersParamsWrapper struct {
	// Just reach endpoint to get info.
	// in:body
	Body struct{}
}

// swagger:route POST /basket/swap foobar-tag idOfswapEndpoint
// swap exchanges currency 1 to currency 2
// responses:
//   200: swapResponse
//   403: swapErrResponse
//   412: swapErrNewResponse

// swagger:response swapResponse
type swapResponseWrapper struct {
	// in:body
	Body router.Transaction
}

// If user doesn't sign in return status 403
// swagger:response swapErrResponse
type swapNewResponseWrapper struct {
	// in:body
	Body struct{}
}

// If user doesn't have enough currency or currency doesn't exist return status 412
// swagger:response swapErrNewResponse
type swapNewNewResponseWrapper struct {
	// in:body
	Body entity.ErrorResponse
}

// swagger:parameters idOfswapEndpoint
type swapParamsWrapper struct {
	// Gives back id of basket with currency
	// in:body
	Body struct {
		bId int `json:"bid"`
	}
}
