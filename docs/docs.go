// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag/v2"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "components": {"schemas":{"dto.Address":{"properties":{"district":{"type":"string"},"province":{"type":"string"},"subdistrict":{"type":"string"},"zipcode":{"type":"string"}},"type":"object"},"dto.CreateTransactionResponseBody":{"properties":{"checkoutUrl":{"type":"string"}},"type":"object"},"dto.DormCreateRequestBody":{"properties":{"address":{"properties":{"district":{"type":"string"},"province":{"type":"string"},"subdistrict":{"type":"string"},"zipcode":{"type":"string"}},"required":["district","province","subdistrict","zipcode"],"type":"object"},"bathrooms":{"minimum":0,"type":"integer"},"bedrooms":{"minimum":0,"type":"integer"},"description":{"type":"string"},"name":{"type":"string"},"price":{"type":"number"},"size":{"type":"number"}},"required":["address","bathrooms","bedrooms","name","price","size"],"type":"object"},"dto.DormImageUploadResponseBody":{"properties":{"url":{"type":"string"}},"type":"object"},"dto.DormResponseBody":{"properties":{"address":{"$ref":"#/components/schemas/dto.Address"},"bathrooms":{"type":"integer"},"bedrooms":{"type":"integer"},"createAt":{"type":"string"},"description":{"type":"string"},"id":{"type":"string"},"imagesUrl":{"items":{"type":"string"},"type":"array","uniqueItems":false},"name":{"type":"string"},"owner":{"$ref":"#/components/schemas/dto.UserResponse"},"price":{"type":"number"},"rating":{"type":"number"},"size":{"type":"number"},"updateAt":{"type":"string"}},"type":"object"},"dto.DormUpdateRequestBody":{"properties":{"address":{"$ref":"#/components/schemas/dto.Address"},"bathrooms":{"minimum":0,"type":"integer"},"bedrooms":{"minimum":0,"type":"integer"},"description":{"type":"string"},"name":{"type":"string"},"price":{"type":"number"},"size":{"type":"number"}},"type":"object"},"dto.ErrorResponse":{"properties":{"error":{"type":"string"}},"type":"object"},"dto.LeasingHistory":{"properties":{"dorm":{"$ref":"#/components/schemas/dto.DormResponseBody"},"end":{"type":"string"},"id":{"type":"string"},"lessee":{"$ref":"#/components/schemas/dto.UserResponse"},"orders":{"items":{"$ref":"#/components/schemas/dto.OrderResponseBody"},"type":"array","uniqueItems":false},"price":{"type":"number"},"start":{"type":"string"}},"type":"object"},"dto.LoginRequestBody":{"properties":{"email":{"type":"string"},"password":{"type":"string"}},"required":["email","password"],"type":"object"},"dto.OrderRequestBody":{"properties":{"leasingHistoryId":{"type":"string"}},"required":["leasingHistoryId"],"type":"object"},"dto.OrderResponseBody":{"properties":{"id":{"type":"string"},"paidTransaction":{"$ref":"#/components/schemas/dto.TransactionResponse"},"price":{"type":"integer"},"type":{"type":"string"}},"type":"object"},"dto.OwnershipProofResponseBody":{"properties":{"adminId":{"type":"string"},"dormId":{"type":"string"},"status":{"$ref":"#/components/schemas/dto.OwnershipProofStatus"},"url":{"type":"string"}},"type":"object"},"dto.OwnershipProofStatus":{"type":"string","x-enum-varnames":["Pending","Approved","Rejected"]},"dto.Pagination":{"properties":{"current_page":{"type":"integer"},"last_page":{"type":"integer"},"limit":{"type":"integer"},"total":{"type":"integer"}},"type":"object"},"dto.PaginationResponse-dto_DormResponseBody":{"properties":{"data":{"items":{"$ref":"#/components/schemas/dto.DormResponseBody"},"type":"array","uniqueItems":false},"pagination":{"$ref":"#/components/schemas/dto.Pagination"}},"type":"object"},"dto.PaginationResponse-dto_LeasingHistory":{"properties":{"data":{"items":{"$ref":"#/components/schemas/dto.LeasingHistory"},"type":"array","uniqueItems":false},"pagination":{"$ref":"#/components/schemas/dto.Pagination"}},"type":"object"},"dto.PaginationResponse-dto_OrderResponseBody":{"properties":{"data":{"items":{"$ref":"#/components/schemas/dto.OrderResponseBody"},"type":"array","uniqueItems":false},"pagination":{"$ref":"#/components/schemas/dto.Pagination"}},"type":"object"},"dto.RefreshTokenRequestBody":{"properties":{"refreshToken":{"type":"string"}},"required":["refreshToken"],"type":"object"},"dto.RegisterRequestBody":{"properties":{"email":{"type":"string"},"password":{"type":"string"},"username":{"type":"string"}},"required":["email","password","username"],"type":"object"},"dto.ResetPasswordCreateRequestBody":{"properties":{"email":{"type":"string"}},"required":["email"],"type":"object"},"dto.ResetPasswordRequestBody":{"properties":{"password":{"type":"string"},"token":{"type":"string"}},"required":["password","token"],"type":"object"},"dto.SuccessResponse-dto_CreateTransactionResponseBody":{"properties":{"data":{"$ref":"#/components/schemas/dto.CreateTransactionResponseBody"}},"type":"object"},"dto.SuccessResponse-dto_DormImageUploadResponseBody":{"properties":{"data":{"$ref":"#/components/schemas/dto.DormImageUploadResponseBody"}},"type":"object"},"dto.SuccessResponse-dto_DormResponseBody":{"properties":{"data":{"$ref":"#/components/schemas/dto.DormResponseBody"}},"type":"object"},"dto.SuccessResponse-dto_LeasingHistory":{"properties":{"data":{"$ref":"#/components/schemas/dto.LeasingHistory"}},"type":"object"},"dto.SuccessResponse-dto_OrderResponseBody":{"properties":{"data":{"$ref":"#/components/schemas/dto.OrderResponseBody"}},"type":"object"},"dto.SuccessResponse-dto_OwnershipProofResponseBody":{"properties":{"data":{"$ref":"#/components/schemas/dto.OwnershipProofResponseBody"}},"type":"object"},"dto.SuccessResponse-dto_TokenResponseBody":{"properties":{"data":{"$ref":"#/components/schemas/dto.TokenResponseBody"}},"type":"object"},"dto.SuccessResponse-dto_TokenWithUserInformationResponseBody":{"properties":{"data":{"$ref":"#/components/schemas/dto.TokenWithUserInformationResponseBody"}},"type":"object"},"dto.SuccessResponse-dto_UserResponse":{"properties":{"data":{"$ref":"#/components/schemas/dto.UserResponse"}},"type":"object"},"dto.TokenResponseBody":{"properties":{"accessToken":{"type":"string"},"refreshToken":{"type":"string"}},"type":"object"},"dto.TokenWithUserInformationResponseBody":{"properties":{"accessToken":{"type":"string"},"refreshToken":{"type":"string"},"userInformation":{"$ref":"#/components/schemas/dto.UserResponse"}},"type":"object"},"dto.TransactionRequestBody":{"properties":{"orderID":{"type":"string"}},"type":"object"},"dto.TransactionResponse":{"properties":{"createAt":{"type":"string"},"id":{"type":"string"},"price":{"type":"integer"},"status":{"type":"string"},"updateAt":{"type":"string"}},"type":"object"},"dto.UserInformationRequestBody":{"properties":{"birthDate":{"type":"string"},"firstname":{"minLength":2,"type":"string"},"gender":{"type":"string"},"lastname":{"minLength":2,"type":"string"},"lifestyles":{"items":{"type":"string"},"type":"array","uniqueItems":false},"nationalID":{"type":"string"},"password":{"minLength":8,"type":"string"},"phoneNumber":{"type":"string"},"studentEvidence":{"type":"string"},"username":{"minLength":2,"type":"string"}},"type":"object"},"dto.UserResponse":{"properties":{"birthDate":{"type":"string"},"email":{"type":"string"},"filledPersonalInfo":{"type":"boolean"},"firstname":{"type":"string"},"gender":{"type":"string"},"id":{"type":"string"},"isStudentVerified":{"type":"boolean"},"isVerified":{"type":"boolean"},"lastname":{"type":"string"},"lifestyles":{"items":{"type":"string"},"type":"array","uniqueItems":false},"phoneNumber":{"type":"string"},"role":{"type":"string"},"studentEvidence":{"type":"string"},"username":{"type":"string"}},"type":"object"},"dto.VerifyRequestBody":{"properties":{"token":{"type":"string"}},"required":["token"],"type":"object"}},"securitySchemes":{"Bearer":{"description":"Bearer token authentication","in":"header","name":"Authorization","type":"apiKey"}}},
    "info": {"description":"{{escape .Description}}","title":"{{.Title}}","version":"{{.Version}}"},
    "externalDocs": {"description":"","url":""},
    "paths": {"/auth/login":{"post":{"description":"Login user","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.LoginRequestBody"}}},"description":"user information","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_TokenWithUserInformationResponseBody"}}},"description":"user successfully logged in"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"user not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"system cannot login user"}},"summary":"Login user","tags":["auth"]}},"/auth/refresh":{"post":{"description":"Refresh user","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.RefreshTokenRequestBody"}}},"description":"user information","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_TokenResponseBody"}}},"description":"user successfully Refresh in"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"user not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"system cannot refresh user"}},"summary":"Refresh user","tags":["auth"]}},"/auth/register":{"post":{"description":"Register new user","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.RegisterRequestBody"}}},"description":"user information","required":true},"responses":{"201":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_TokenWithUserInformationResponseBody"}}},"description":"user successfully registered"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"system cannot register user"}},"summary":"Register new user","tags":["auth"]}},"/dorms":{"get":{"description":"Retrieve a list of all dorms filtered by name. If no query are provided, all dorms are returned.","parameters":[{"description":"Dorm name to search","in":"query","name":"name","schema":{"type":"string"}},{"description":"Number of dorms to retrieve (default 10, max 50)","in":"query","name":"limit","schema":{"type":"integer"}},{"description":"Page number to retrieve (default 1)","in":"query","name":"page","schema":{"type":"integer"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.PaginationResponse-dto_DormResponseBody"}}},"description":"All dorms retrieved successfully"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Failed to retrieve dorms"}},"summary":"Get all dorms by name","tags":["dorms"]},"post":{"description":"Add a new room to the database with the given details","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.DormCreateRequestBody"}}},"description":"Dorm information","required":true},"responses":{"201":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_DormResponseBody"}}},"description":"Dorm successfully created"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Your request is invalid"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"403":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"You do not have permission to create a dorm"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Failed to save dorm"}},"security":[{"Bearer":[]}],"summary":"Create a new dorm","tags":["dorms"]}},"/dorms/{id}":{"delete":{"description":"Removes a dorm from the database based on the give ID","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"204":{"description":"Dorm successfully deleted"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"403":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"You do not have permission to delete this dorm"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Dorm not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Failed to delete dorm"}},"security":[{"Bearer":[]}],"summary":"Delete a dorm","tags":["dorms"]},"get":{"description":"Retrieve a specific dorm based on its ID","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_DormResponseBody"}}},"description":"Dorm data successfully retrieved"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Dorm not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Server failed to retrieve dorm"}},"summary":"Get a dorm by ID","tags":["dorms"]},"patch":{"description":"Modifies an existing room's details based on the given ID","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.DormUpdateRequestBody"}}},"description":"Updated Room Data","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_DormResponseBody"}}},"description":"Dorm data updated successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Invalid Request"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"403":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"unauthorized to update this dorm"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Dorm not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Server failed to update dorm"}},"security":[{"Bearer":[]}],"summary":"Update an existing dorm","tags":["dorms"]}},"/dorms/{id}/images":{"post":{"description":"Upload an image for a specific dorm by its ID, by attaching the image as a value for the key field name \"image\", as a multipart form-data","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"requestBody":{"content":{"multipart/form-data":{"schema":{"type":"file"}}},"description":"DormImage","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_DormImageUploadResponseBody"}}},"description":"Successful image upload"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Invalid Request"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"403":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"unauthorized to upload image to dorm"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Dorm not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Server failed to upload dorm image"}},"security":[{"Bearer":[]}],"summary":"Upload an image for a dorm","tags":["dorms"]}},"/history/bydorm/{id}":{"get":{"description":"Retrieve a list of all leasing history by userid","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}},{"description":"Number of dorms to retrieve (default 10, max 50)","in":"query","name":"limit","schema":{"type":"integer"}},{"description":"Page number to retrieve (default 1)","in":"query","name":"page","schema":{"type":"integer"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.PaginationResponse-dto_LeasingHistory"}}},"description":"Retrive history successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format or limit parameter is incorrect or page parameter is incorrect or page exceeded"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"leasing history not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Can not parse UUID"}},"security":[{"Bearer":[]}],"summary":"Get all leasing history by userid","tags":["history"]}},"/history/me":{"get":{"description":"Retrieve a list of all leasing history by userid","parameters":[{"description":"Number of dorms to retrieve (default 10, max 50)","in":"query","name":"limit","schema":{"type":"integer"}},{"description":"Page number to retrieve (default 1)","in":"query","name":"page","schema":{"type":"integer"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.PaginationResponse-dto_LeasingHistory"}}},"description":"Retrive history successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format or limit parameter is incorrect or page parameter is incorrect or page exceeded"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"leasing history not found"}},"security":[{"Bearer":[]}],"summary":"Get all leasing history by userid","tags":["history"]}},"/history/{id}":{"delete":{"description":"Delete a leasing history in the database","parameters":[{"description":"LeasingHistoryId","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"204":{"description":"No Content"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"leasing history not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Can not parse UUID or Failed to delete leasing history"}},"security":[{"Bearer":[]}],"summary":"Delete a leasing history","tags":["history"]},"patch":{"description":"Delete a leasing history in the database","parameters":[{"description":"LeasingHistoryId","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"204":{"description":"Set end timestamp successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"leasing history not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Can not parse UUID or Failed to update leasing history"}},"security":[{"Bearer":[]}],"summary":"Delete a leasing history","tags":["history"]},"post":{"description":"Add a new leasing history to the database","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"201":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_LeasingHistory"}}},"description":"Dorm successfully created"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Dorm not found or leasing history not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Can not parse UUID or failed to save leasing history to database"}},"security":[{"Bearer":[]}],"summary":"Create a new leasing history","tags":["history"]}},"/order":{"post":{"description":"Create an order","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.OrderRequestBody"}}},"description":"Order request body","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_OrderResponseBody"}}},"description":"Order created successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"leasing history not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"cannot parse uuid or cannot delete user"}},"security":[{"Bearer":[]}],"summary":"Create an order","tags":["order"]}},"/order/unpaid/me":{"get":{"description":"Get my unpaid orders by ID","parameters":[{"description":"Number of history to be retrieved","in":"query","name":"limit","schema":{"type":"integer"}},{"description":"Page to retrieved","in":"query","name":"page","schema":{"type":"integer"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.PaginationResponse-dto_OrderResponseBody"}}},"description":"Unpaid orders retrieved successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"leasing history not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"cannot parse uuid or cannot delete user"}},"security":[{"Bearer":[]}],"summary":"Get my unpaid orders by ID","tags":["order"]}},"/order/unpaid/{userID}":{"get":{"description":"Get unpaid orders by User ID","parameters":[{"description":"User ID","in":"path","name":"userID","required":true,"schema":{"type":"string"}},{"description":"Number of history to be retrieved","in":"query","name":"limit","schema":{"type":"integer"}},{"description":"Page to retrieved","in":"query","name":"page","schema":{"type":"integer"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.PaginationResponse-dto_OrderResponseBody"}}},"description":"Order retrieved successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"leasing history not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"cannot parse uuid or cannot delete user"}},"security":[{"Bearer":[]}],"summary":"Get unpaid orders by User ID","tags":["order"]}},"/order/{id}":{"get":{"description":"Get an order by ID","parameters":[{"description":"Order ID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_OrderResponseBody"}}},"description":"Order retrieved successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"leasing history not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"cannot parse uuid or cannot delete user"}},"security":[{"Bearer":[]}],"summary":"Get an order by ID","tags":["order"]}},"/ownership/{id}":{"delete":{"description":"Delete an ownership proof file","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"204":{"description":"Ownership proof successfully deleted"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Ownership file not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Failed to delete Ownership file"}},"security":[{"Bearer":[]}],"summary":"Delete ownership proof","tags":["ownership"]},"get":{"description":"Retrieve ownership proof for a specific dorm","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_OwnershipProofResponseBody"}}},"description":"Ownership proof retrieved successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Ownership file not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Internal server error"}},"security":[{"Bearer":[]}],"summary":"Get ownership proof by Dorm ID","tags":["ownership"]}},"/ownership/{id}/approve":{"post":{"description":"Approve a submitted ownership proof for a dorm","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_OwnershipProofResponseBody"}}},"description":"Ownership proof approved"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Ownership file not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Internal server error"}},"security":[{"Bearer":[]}],"summary":"Approve ownership proof","tags":["ownership"]}},"/ownership/{id}/reject":{"post":{"description":"Reject a submitted ownership proof for a dorm","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_OwnershipProofResponseBody"}}},"description":"Ownership proof rejected"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Ownership file not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Internal server error"}},"security":[{"Bearer":[]}],"summary":"Reject ownership proof","tags":["ownership"]}},"/ownership/{id}/upload":{"post":{"description":"Upload a new file as ownership proof for a dorm","parameters":[{"description":"DormID","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"requestBody":{"content":{"multipart/form-data":{"schema":{"type":"file"}}},"description":"file","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_OwnershipProofResponseBody"}}},"description":"Ownership proof created"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Incorrect UUID format"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Ownershop proof not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"Server failed to upload file"}},"security":[{"Bearer":[]}],"summary":"Upload new ownership proof","tags":["ownership"]}},"/transaction":{"post":{"description":"Create a transaction","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.TransactionRequestBody"}}},"description":"Transaction request body","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_CreateTransactionResponseBody"}}},"description":"Transaction created successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"404":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"leasing history not found"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"cannot parse uuid or cannot delete user"}},"security":[{"Bearer":[]}],"summary":"Create a transaction","tags":["transaction"]}},"/user":{"patch":{"description":"Update user information","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.UserInformationRequestBody"}}},"description":"user information","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_UserResponse"}}},"description":"user successfully updated account information"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"system cannot update your account information"}},"security":[{"Bearer":[]}],"summary":"Update user information","tags":["user"]}},"/user/":{"delete":{"description":"Delete a user account","requestBody":{"content":{"application/json":{"schema":{"type":"object"}}}},"responses":{"204":{"description":"account successfully deleted"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"cannot parse uuid or cannot delete user"}},"security":[{"Bearer":[]}],"summary":"Delete a user account","tags":["user"]}},"/user/me":{"get":{"description":"Get user information","responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_UserResponse"}}},"description":"get user information successfully"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"system cannot get user information"}},"security":[{"Bearer":[]}],"summary":"Get user information","tags":["user"]}},"/user/newpassword":{"post":{"description":"Reset password","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ResetPasswordRequestBody"}}},"description":"token","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_TokenWithUserInformationResponseBody"}}},"description":"password reset successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"system cannot reset password"}},"summary":"Reset password","tags":["user"]}},"/user/resetpassword":{"post":{"description":"Resend verification email","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ResetPasswordCreateRequestBody"}}},"description":"token","required":true},"responses":{"204":{"description":"email is sent to user successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"system cannot resend verification email"}},"summary":"Resend verification email","tags":["user"]}},"/user/verify":{"post":{"description":"Verify email","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.VerifyRequestBody"}}},"description":"token","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_TokenWithUserInformationResponseBody"}}},"description":"email is verified successfully"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is invalid"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"system cannot verify your email"}},"summary":"Verify email","tags":["user"]}},"/user/{id}":{"get":{"description":"Get User By ID","parameters":[{"description":"user id","in":"path","name":"id","required":true,"schema":{"type":"string"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-dto_UserResponse"}}},"description":"get user information successfully"},"401":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"your request is unauthorized"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.ErrorResponse"}}},"description":"system cannot get user information"}},"security":[{"Bearer":[]}],"summary":"GetUserByID","tags":["user"]}}},
    "openapi": "3.1.0"
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
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
