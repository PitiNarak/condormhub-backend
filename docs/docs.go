// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag/v2"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},"swagger":"2.0","info":{"description":"{{escape .Description}}","title":"{{.Title}}","contact":{},"version":"{{.Version}}"},"host":"{{.Host}}","basePath":"{{.BasePath}}","paths":{"/auth/login":{"post":{"description":"Login user","consumes":["application/json"],"produces":["application/json"],"tags":["auth"],"summary":"Login user","parameters":[{"description":"user information","name":"user","in":"body","required":true,"schema":{"$ref":"#/definitions/dto.LoginRequestBody"}}],"responses":{"200":{"description":"user successfully logged in","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/dto.TokenWithUserInformationResponseBody"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"user not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"system cannot login user","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/auth/refresh":{"post":{"description":"Refresh user","consumes":["application/json"],"produces":["application/json"],"tags":["auth"],"summary":"Refresh user","parameters":[{"description":"user information","name":"user","in":"body","required":true,"schema":{"$ref":"#/definitions/dto.RefreshTokenRequestBody"}}],"responses":{"200":{"description":"user successfully Refresh in","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/dto.TokenResponseBody"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"user not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"system cannot refresh user","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/auth/register":{"post":{"description":"Register new user","consumes":["application/json"],"produces":["application/json"],"tags":["auth"],"summary":"Register new user","parameters":[{"description":"user information","name":"user","in":"body","required":true,"schema":{"$ref":"#/definitions/dto.RegisterRequestBody"}}],"responses":{"201":{"description":"user successfully registered","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/dto.TokenWithUserInformationResponseBody"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"system cannot register user","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/dorms":{"get":{"description":"Retrieve a list of all dorms","produces":["application/json"],"tags":["dorms"],"summary":"Get all dorms","responses":{"200":{"description":"All dorms retrieved successfully","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"array","items":{"$ref":"#/definitions/domain.Dorm"}},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"Failed to retrieve dorms","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}},"post":{"security":[{"Bearer":[]}],"description":"Add a new room to the database with the given details","consumes":["application/json"],"produces":["application/json"],"tags":["dorms"],"summary":"Create a new dorm","parameters":[{"description":"Dorm information","name":"dorm","in":"body","required":true,"schema":{"$ref":"#/definitions/dto.DormRequestBody"}}],"responses":{"201":{"description":"Dorm successfully created","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/domain.Dorm"},"pagination":{"type":"object"}}}]}},"400":{"description":"Your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"Failed to save dorm","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/dorms/{id}":{"get":{"description":"Retrieve a specific dorm based on its ID","produces":["application/json"],"tags":["dorms"],"summary":"Get a dorm by ID","parameters":[{"type":"string","description":"DormID","name":"id","in":"path","required":true}],"responses":{"200":{"description":"Dorm data successfully retrieved","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/domain.Dorm"},"pagination":{"type":"object"}}}]}},"400":{"description":"Incorrect UUID format","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"Dorm not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"Server failed to retrieve dorm","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}},"delete":{"security":[{"Bearer":[]}],"description":"Removes a dorm from the database based on the give ID","produces":["application/json"],"tags":["dorms"],"summary":"Delete a dorm","parameters":[{"type":"string","description":"DormID","name":"id","in":"path","required":true}],"responses":{"200":{"description":"Dorm successfully deleted","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"400":{"description":"Incorrect UUID format","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"Dorm not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"Failed to delete dorm","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}},"patch":{"security":[{"Bearer":[]}],"description":"Modifies an existing room's details based on the given ID","consumes":["application/json"],"produces":["application/json"],"tags":["dorms"],"summary":"Update an existing dorm","parameters":[{"type":"string","description":"DormID","name":"id","in":"path","required":true},{"description":"Updated Room Data","name":"dorm","in":"body","required":true,"schema":{"$ref":"#/definitions/dto.DormRequestBody"}}],"responses":{"200":{"description":"Dorm data updated successfully","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/domain.Dorm"},"pagination":{"type":"object"}}}]}},"400":{"description":"Invalid Request","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"Dorm not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"Server failed to update dorm","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/history/bydorm/{id}":{"get":{"security":[{"Bearer":[]}],"description":"Retrieve a list of all leasing history by userid","produces":["application/json"],"tags":["history"],"summary":"Get all leasing history by userid","parameters":[{"type":"string","description":"DormID","name":"id","in":"path","required":true},{"type":"string","description":"Number of history to be retirved","name":"limit","in":"query","required":true},{"type":"string","description":"Page to retrive","name":"page","in":"query","required":true}],"responses":{"200":{"description":"Retrive history successfully","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"array","items":{"$ref":"#/definitions/domain.LeasingHistory"}},"pagination":{"$ref":"#/definitions/dto.PaginationResponseBody"}}}]}},"400":{"description":"Incorrect UUID format or limit parameter is incorrect or page parameter is incorrect or page exceeded","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"leasing history not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"Can not parse UUID","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/history/me":{"get":{"security":[{"Bearer":[]}],"description":"Retrieve a list of all leasing history by userid","produces":["application/json"],"tags":["history"],"summary":"Get all leasing history by userid","parameters":[{"type":"string","description":"Number of history to be retirved","name":"limit","in":"query","required":true},{"type":"string","description":"Page to retrive","name":"page","in":"query","required":true}],"responses":{"200":{"description":"Retrive history successfully","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"array","items":{"$ref":"#/definitions/domain.LeasingHistory"}},"pagination":{"$ref":"#/definitions/dto.PaginationResponseBody"}}}]}},"400":{"description":"Incorrect UUID format or limit parameter is incorrect or page parameter is incorrect or page exceeded","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"leasing history not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/history/{id}":{"post":{"security":[{"Bearer":[]}],"description":"Add a new leasing history to the database","produces":["application/json"],"tags":["history"],"summary":"Create a new leasing history","parameters":[{"type":"string","description":"DormID","name":"id","in":"path","required":true}],"responses":{"201":{"description":"Dorm successfully created","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/domain.LeasingHistory"},"pagination":{"type":"object"}}}]}},"400":{"description":"Incorrect UUID format","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"Dorm not found or leasing history not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"Can not parse UUID or failed to save leasing history to database","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}},"delete":{"security":[{"Bearer":[]}],"description":"Delete a leasing history in the database","produces":["application/json"],"tags":["history"],"summary":"Delete a leasing history","parameters":[{"type":"string","description":"LeasingHistoryId","name":"id","in":"path","required":true}],"responses":{"200":{"description":"Delete successfully","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"400":{"description":"Incorrect UUID format","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"leasing history not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"Can not parse UUID or Failed to delete leasing history","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}},"patch":{"security":[{"Bearer":[]}],"description":"Delete a leasing history in the database","produces":["application/json"],"tags":["history"],"summary":"Delete a leasing history","parameters":[{"type":"string","description":"LeasingHistoryId","name":"id","in":"path","required":true}],"responses":{"200":{"description":"Set end timestamp successfully","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"400":{"description":"Incorrect UUID format","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"leasing history not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"Can not parse UUID or Failed to update leasing history","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/order":{"post":{"security":[{"Bearer":[]}],"description":"Create an order","consumes":["application/json"],"produces":["application/json"],"tags":["order"],"summary":"Create an order","parameters":[{"description":"Order request body","name":"body","in":"body","required":true,"schema":{"$ref":"#/definitions/dto.OrderRequestBody"}}],"responses":{"200":{"description":"account successfully deleted","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/dto.OrderResponseBody"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"leasing history not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"cannot parse uuid or cannot delete user","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/order/unpaid/me":{"get":{"security":[{"Bearer":[]}],"description":"Get my unpaid orders by ID","consumes":["application/json"],"produces":["application/json"],"tags":["order"],"summary":"Get my unpaid orders by ID","parameters":[{"type":"string","description":"Number of history to be retrieved","name":"limit","in":"query","required":true},{"type":"string","description":"Page to retrieved","name":"page","in":"query","required":true}],"responses":{"200":{"description":"account successfully deleted","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/dto.OrderResponseBody"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"leasing history not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"cannot parse uuid or cannot delete user","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/order/unpaid/{userID}":{"get":{"security":[{"Bearer":[]}],"description":"Get unpaid orders by ID","consumes":["application/json"],"produces":["application/json"],"tags":["order"],"summary":"Get unpaid orders by ID","parameters":[{"type":"string","description":"User ID","name":"userID","in":"path","required":true},{"type":"string","description":"Number of history to be retrieved","name":"limit","in":"query","required":true},{"type":"string","description":"Page to retrieved","name":"page","in":"query","required":true}],"responses":{"200":{"description":"account successfully deleted","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/dto.OrderResponseBody"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"leasing history not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"cannot parse uuid or cannot delete user","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/order/{id}":{"get":{"security":[{"Bearer":[]}],"description":"Get an order by ID","consumes":["application/json"],"produces":["application/json"],"tags":["order"],"summary":"Get an order by ID","parameters":[{"type":"string","description":"Order ID","name":"id","in":"path","required":true}],"responses":{"200":{"description":"account successfully deleted","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/dto.OrderResponseBody"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"404":{"description":"leasing history not found","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"cannot parse uuid or cannot delete user","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/user":{"patch":{"security":[{"Bearer":[]}],"description":"Update user information","consumes":["application/json"],"produces":["application/json"],"tags":["user"],"summary":"Update user information","parameters":[{"description":"user information","name":"user","in":"body","required":true,"schema":{"$ref":"#/definitions/dto.UserInformationRequestBody"}}],"responses":{"200":{"description":"user successfully updated account information","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/domain.User"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"system cannot update your account information","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/user/":{"delete":{"security":[{"Bearer":[]}],"description":"Delete a user account","consumes":["application/json"],"produces":["application/json"],"tags":["user"],"summary":"Delete a user account","responses":{"200":{"description":"account successfully deleted","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"cannot parse uuid or cannot delete user","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/user/me":{"get":{"security":[{"Bearer":[]}],"description":"Get user information","consumes":["application/json"],"produces":["application/json"],"tags":["user"],"summary":"Get user information","responses":{"200":{"description":"get user information successfully","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/domain.User"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"system cannot get user information","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/user/newpassword":{"post":{"description":"Reset password","consumes":["application/json"],"produces":["application/json"],"tags":["user"],"summary":"Reset password","parameters":[{"description":"token","name":"user","in":"body","required":true,"schema":{"$ref":"#/definitions/dto.ResetPasswordRequestBody"}}],"responses":{"200":{"description":"password reset successfully","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/dto.TokenWithUserInformationResponseBody"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"system cannot reset password","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/user/resetpassword":{"post":{"description":"Resend verification email","consumes":["application/json"],"produces":["application/json"],"tags":["user"],"summary":"Resend verification email","parameters":[{"description":"token","name":"user","in":"body","required":true,"schema":{"$ref":"#/definitions/dto.ResetPasswordCreateRequestBody"}}],"responses":{"200":{"description":"email is sent to user successfully","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"system cannot resend verification email","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}},"/user/verify":{"post":{"description":"Verify email","consumes":["application/json"],"produces":["application/json"],"tags":["user"],"summary":"Verify email","parameters":[{"description":"token","name":"user","in":"body","required":true,"schema":{"$ref":"#/definitions/dto.VerifyRequestBody"}}],"responses":{"200":{"description":"email is verified successfully","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"$ref":"#/definitions/dto.TokenWithUserInformationResponseBody"},"pagination":{"type":"object"}}}]}},"400":{"description":"your request is invalid","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"401":{"description":"your request is unauthorized","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}},"500":{"description":"system cannot verify your email","schema":{"allOf":[{"$ref":"#/definitions/httpResponse.HttpResponse"},{"type":"object","properties":{"data":{"type":"object"},"pagination":{"type":"object"}}}]}}}}}},"definitions":{"domain.Address":{"type":"object","required":["district","province","subdistrict","zipcode"],"properties":{"district":{"type":"string"},"province":{"type":"string"},"subdistrict":{"type":"string"},"zipcode":{"type":"string"}}},"domain.CheckoutStatus":{"type":"string","enum":["open","complete","expired"],"x-enum-varnames":["StatusOpen","StatusComplete","StatusExpired"]},"domain.Dorm":{"type":"object","required":["address","bathrooms","bedrooms","name","ownerId","price","size"],"properties":{"address":{"$ref":"#/definitions/domain.Address"},"bathrooms":{"type":"integer","minimum":0},"bedrooms":{"type":"integer","minimum":0},"createAt":{"type":"string"},"description":{"type":"string"},"id":{"type":"string"},"name":{"type":"string"},"owner":{"$ref":"#/definitions/domain.User"},"ownerId":{"type":"string"},"price":{"type":"number"},"rating":{"type":"number","maximum":5,"minimum":0},"size":{"type":"number"},"updateAt":{"type":"string"}}},"domain.LeasingHistory":{"type":"object","properties":{"dorm":{"$ref":"#/definitions/domain.Dorm"},"dorm_id":{"type":"string"},"end":{"type":"string"},"id":{"type":"string"},"lessee":{"$ref":"#/definitions/domain.User"},"lessee_id":{"type":"string"},"orders":{"type":"array","items":{"$ref":"#/definitions/domain.Order"}},"start":{"type":"string"}}},"domain.Lifestyle":{"type":"string","enum":["Active","Creative","Social","Relaxed","Football","Basketball","Tennis","Swimming","Running","Cycling","Badminton","Yoga","Gym \u0026 Fitness","Music","Dancing","Photography","Painting","Gaming","Reading","Writing","DIY \u0026 Crafting","Cooking","Extrovert","Introvert","Night Owl","Early Bird","Traveler","Backpacker","Nature Lover","Camping","Beach Lover","Dog Lover","Cat Lover","Freelancer","Entrepreneur","Office Worker","Remote Worker","Student","Self-Employed"],"x-enum-varnames":["Active","Creative","Social","Relaxed","Football","Basketball","Tennis","Swimming","Running","Cycling","Badminton","Yoga","GymAndFitness","Music","Dancing","Photography","Painting","Gaming","Reading","Writing","DIYAndCrafting","Cooking","Extrovert","Introvert","NightOwl","EarlyBird","Traveler","Backpacker","NatureLover","Camping","BeachLover","DogLover","CatLover","Freelancer","Entrepreneur","OfficeWorker","RemoteWorker","Student","SelfEmployed"]},"domain.Order":{"type":"object","properties":{"createAt":{"type":"string"},"id":{"type":"string"},"leasingHistory":{"$ref":"#/definitions/domain.LeasingHistory"},"leasingHistoryID":{"type":"string"},"paidTransaction":{"$ref":"#/definitions/domain.Transaction"},"paidTransactionID":{"type":"string"},"price":{"type":"integer"},"transactions":{"type":"array","items":{"$ref":"#/definitions/domain.Transaction"}},"type":{"$ref":"#/definitions/domain.OrderType"},"updateAt":{"type":"string"}}},"domain.OrderType":{"type":"string","enum":["insurance","monthly_bill"],"x-enum-varnames":["InsuranceOrderType","MonthlyBillOrderType"]},"domain.Role":{"type":"string","enum":["ADMIN","LESSEE","LESSOR"],"x-enum-varnames":["AdminRole","LesseeRole","LessorRole"]},"domain.Transaction":{"type":"object","properties":{"createAt":{"type":"string"},"id":{"type":"string"},"order":{"$ref":"#/definitions/domain.Order"},"orderID":{"type":"string"},"price":{"type":"integer"},"status":{"$ref":"#/definitions/domain.CheckoutStatus"},"updateAt":{"type":"string"}}},"domain.User":{"type":"object","required":["email","username"],"properties":{"birthDate":{"type":"string"},"createAt":{"type":"string"},"email":{"type":"string"},"filledPersonalInfo":{"type":"boolean"},"firstname":{"type":"string"},"gender":{"type":"string"},"id":{"type":"string"},"isStudentVerified":{"type":"boolean"},"isVerified":{"type":"boolean"},"lastname":{"type":"string"},"lifestyles":{"type":"array","items":{"$ref":"#/definitions/domain.Lifestyle"}},"nationalID":{"type":"string"},"role":{"$ref":"#/definitions/domain.Role"},"studentEvidence":{"description":"studentEvidence","type":"string"},"updateAt":{"type":"string"},"username":{"type":"string"}}},"dto.DormRequestBody":{"type":"object","required":["address","bathrooms","bedrooms","name","price","size"],"properties":{"address":{"type":"object","required":["district","province","subdistrict","zipcode"],"properties":{"district":{"type":"string"},"province":{"type":"string"},"subdistrict":{"type":"string"},"zipcode":{"type":"string"}}},"bathrooms":{"type":"integer","minimum":0},"bedrooms":{"type":"integer","minimum":0},"description":{"type":"string"},"name":{"type":"string"},"price":{"type":"number"},"size":{"type":"number"}}},"dto.LoginRequestBody":{"type":"object","required":["email","password"],"properties":{"email":{"type":"string"},"password":{"type":"string"}}},"dto.OrderRequestBody":{"type":"object","required":["leasingHistoryId"],"properties":{"leasingHistoryId":{"type":"string"}}},"dto.OrderResponseBody":{"type":"object","properties":{"id":{"type":"string"},"paidTransaction":{"$ref":"#/definitions/domain.Transaction"},"price":{"type":"integer"},"transactions":{"type":"array","items":{"$ref":"#/definitions/domain.Transaction"}},"type":{"$ref":"#/definitions/domain.OrderType"}}},"dto.PaginationResponseBody":{"type":"object","properties":{"currentPage":{"type":"integer"},"lastPage":{"type":"integer"},"limit":{"type":"integer"},"total":{"type":"integer"}}},"dto.RefreshTokenRequestBody":{"type":"object","required":["refreshToken"],"properties":{"refreshToken":{"type":"string"}}},"dto.RegisterRequestBody":{"type":"object","required":["email","password","username"],"properties":{"email":{"type":"string"},"password":{"type":"string"},"username":{"type":"string"}}},"dto.ResetPasswordCreateRequestBody":{"type":"object","required":["email"],"properties":{"email":{"type":"string"}}},"dto.ResetPasswordRequestBody":{"type":"object","required":["password","token"],"properties":{"password":{"type":"string"},"token":{"type":"string"}}},"dto.TokenResponseBody":{"type":"object","properties":{"accessToken":{"type":"string"},"refreshToken":{"type":"string"}}},"dto.TokenWithUserInformationResponseBody":{"type":"object","properties":{"accessToken":{"type":"string"},"refreshToken":{"type":"string"},"userInformation":{"$ref":"#/definitions/domain.User"}}},"dto.UserInformationRequestBody":{"type":"object","properties":{"birthDate":{"type":"string"},"firstname":{"type":"string"},"gender":{"type":"string"},"lastname":{"type":"string"},"lifestyles":{"type":"array","items":{"$ref":"#/definitions/domain.Lifestyle"}},"nationalID":{"type":"string"},"password":{"type":"string","minLength":8},"studentEvidence":{"type":"string"},"username":{"type":"string"}}},"dto.VerifyRequestBody":{"type":"object","required":["token"],"properties":{"token":{"type":"string"}}},"httpResponse.HttpResponse":{"type":"object","properties":{"data":{},"message":{"type":"string"},"pagination":{},"success":{"type":"boolean"}}}},"securityDefinitions":{"Bearer":{"description":"Bearer token authentication","type":"apiKey","name":"Authorization","in":"header"}}}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Condormhub API",
	Description:      "This is the API for the Condormhub project.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
