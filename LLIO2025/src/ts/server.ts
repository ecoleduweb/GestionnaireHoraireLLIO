const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

const getHeaders = (token = localStorage.getItem("token")) => ({
    "Content-Type": "application/json",
    Authorization: token || "",
})

const handleResponse = async <T>(response: Response, redirectToLoginOn401 = true, redirectTo500 = true): Promise<T | undefined> => {
    if (!response.ok) {
        const statusCodeRedirects = {
            500: "/500",
            401: "/",
        }

        const shouldRedirect =
            (response.status === 401 && redirectToLoginOn401) ||
            (response.status === 500 && redirectTo500);

        if (shouldRedirect && statusCodeRedirects[response.status]) {
            window.location.href = statusCodeRedirects[response.status];
            return undefined; // type hint for TS
        }

        if (response.status === 403) throw new Error("Erreur - forbidden")

        if (response.status === 404) return undefined as T

        try {
            // Vérifier si la réponse contient du JSON
            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                const errorData = await response.json();

                // Si on retourne un code d'erreur pour distinguer l'erreur, on l'inclus.
                if (errorData.code) {
                  const err = new Error(errorData.error);
                  (err as any).code = errorData.code;
                  throw err;
                }
                // Si l'API retourne un message d'erreur, l'utiliser
                if (errorData.error) {
                    throw new Error(errorData.error);
                }
                
                // Si c'est un autre format JSON, essayer d'autres propriétés communes
                if (errorData.message) {
                    throw new Error(errorData.message);
                }
            }
        } catch (jsonError) {
            // Si c'est notre erreur personnalisée avec le message de l'API, la relancer
            if (jsonError instanceof Error && jsonError.message !== `Error: ${response.status} - ${response.statusText}`) {
                throw jsonError;
            }
            // Sinon, continuer avec l'erreur générique ci-dessous
        }

        throw new Error(`Error: ${response.status} - ${response.statusText}`)
    }
    return response.json() as Promise<T>
}

const request = async <T>(
    method: string,
    url: string,
    body?: any,
    redirectToLoginOn401 = true,
    redirectTo500 = true
): Promise<T> =>{
    try {
        const options: RequestInit = {
            method,
            credentials: 'include',
            headers: getHeaders(),
        }

        if (body) {
            options.body = JSON.stringify(body)
        }

        const response = await fetch(`${VITE_API_BASE_URL}${url}`, options)
        const data = await handleResponse<T>(response, redirectToLoginOn401, redirectTo500)
        return data as T
    } catch (error) {
        console.error(`Error ${method}ing:`, error)
        throw error
    }
}

export const GET = <T>(url: string, redirectToLoginOn401?: boolean): Promise<T> =>
    request<T>("GET", url, undefined, redirectToLoginOn401)

export const POST = <T, R>(url: string, body: T, redirectToLoginOn401?: boolean): Promise<R> =>
    request<R>("POST", url, body, redirectToLoginOn401)

export const DELETE = (url: string, redirectToLoginOn401?: boolean, redirectTo500?: boolean): Promise<void> =>
    request<void>("DELETE", url, undefined, redirectToLoginOn401, redirectTo500)

export const PUT = <T, R>(url: string, body: T, redirectToLoginOn401?: boolean): Promise<R> =>
    request<R>("PUT", url, body, redirectToLoginOn401)

export const PATCH = <T>(url: string, body: T): Promise<void> =>
    request<void>("PATCH", url, body)