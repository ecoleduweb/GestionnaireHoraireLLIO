import { MockConfig } from "../types";

export const userMocks = {
    userMeSuccess: {
        url: '/user/me',
        method: 'GET',
        response: {
                status: 200,
                json:
                {
                    "firstName": "Usager",
                    "lastName": "test",
                    "role": 2,
                }
        }
    },
    usersSuccess: {
        url: '/users',
        method: 'GET',
        response: {
                status: 200,
                json:
                [
                    {
                        "id": 1,
                        "firstName": "Usager",
                        "lastName": "test",
                        "role": 2,
                        "email": "usager@email.com"
                    },
                    {
                        "id": 2,
                        "firstName": "JérémieTest",
                        "lastName": "Lapointe",
                        "role": 1,
                        "email": "jayboss@email.com"
                    },
                    {
                        "id": 3,
                        "firstName": "Charle-ÉtienneTest",
                        "lastName": "Soucy",
                        "role": 0,
                        "email" : "wong@email.com"
                    },
                    {
                        "id": 4,
                        "firstName": "Quentin",
                        "lastName": "Lecourt",
                        "role": 2,
                        "email": "Ariales@email.com"

                    }
                ]
        }
    },
    updateUserRoleSuccess: {
        url: '/user/*/role',
        method: 'PATCH',
        response: {
                status: 200,
                json:
                {
                    "id": 2,
                    "firstName": "JérémieTest",
                    "lastName": "Lapointe",
                    "email": "jeremietest.lapointe@llio.com",
                    "role": 2
                }
        }
    },
    updateUserRoleError: {
        url: '/user/*/role',
        method: 'PATCH',
        response: {
                status: 400,
                json:
                {
                    "error": "Erreur lors de la mise à jour du rôle"
                }
        }
    },
    getAllManagersSuccess: {
        url: '/users?role=1&role=2',
        method: 'GET',
        response: {
                status: 200,
                json:
                [
                    {
                        "id": 2,
                        "firstName": "JérémieTest",
                        "lastName": "Lapointe",
                        "role": 1
                    },
                    {
                        "id": 3,
                        "firstName": "Charle-ÉtienneTest",
                        "lastName": "Soucy",
                        "role": 2
                    }
                ]
        }
    },
    getAllManagersError: {
        url: '/users/managers',
        method: 'GET',
        response: {
                status: 400,
                json:
                {
                    "error": "Erreur lors de la récupération des managers"
                }
        }
    },
    logOutSuccess: {
        url: '/logout',
        method: 'POST',
        response: {
                status: 200,
                json: {"message":"Déconnexion réussie"}
        }
    },
    deleteUserSuccess: {
        url: '/user/*',
        method: 'DELETE',
        response: {
                status: 200,
                json: {"message":"Utilisateur supprimé avec succès"}
        }
    },
    deleteUserError: {
        url: '/user/*',
        method: 'DELETE',
        response: {
                status: 403,
                json: {"error":"L'utilisateur a une ou des activités associées, suppression impossible"}
        }
    },
   
   saveTimeBankConfigSuccess: {
      url: '/user/time-bank/config',
        method: 'POST',
        response: {
            status: 200,
            json: {
                "startDate": "2025-03-01",
                "hoursPerWeek": 40,
                "offset": 0
            }
        }
    },
    getTimeBankConfigSuccess: {
      url: '/user/time-bank/config',
        method: 'GET',
        response: {
            status: 200,
            json: {
                "startDate": "2025-03-01",
                "hoursPerWeek": 40,
                "offset": 0
            }
        }
    },
    getTimeBankSuccess: {
        url: '/user/time-bank',
        method: 'GET',
        response: {
            status: 200,
            json: {
                "isConfigured": true,
                "timeInBank": 40
            }
        }
    },

} satisfies Record<string, MockConfig>;