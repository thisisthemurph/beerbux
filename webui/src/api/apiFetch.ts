const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

type ValidationErrorResponse = {
	errors: Record<string, string>;
};

export class ValidationError extends Error {
	constructor(public validationErrors: ValidationErrorResponse) {
		super("Validation error");
		this.name = "ValidationError";
	}
}

export async function apiFetch<T>(
	url: string,
	options?: RequestInit,
): Promise<T> {
	try {
		const response = await fetch(`${API_BASE_URL}${url}`, {
			...options,
			credentials: "include",
			headers: {
				"Content-Type": "application/json",
				...(options?.headers ?? {}),
			},
		});

		const data = await response.json().catch(() => {
			throw new Error("Invalid JSON response");
		});

		if (!response.ok) {
			throw isValidationErrorResponse(data)
				? new ValidationError(data)
				: new Error(data.error ?? "An error occurred");
		}

		return data as T;
	} catch (error) {
		if (error instanceof ValidationError) {
			throw error; // Rethrow validation errors
		}
		throw new Error(
			error instanceof Error ? error.message : "An error occurred",
		);
	}
}

/*
 * Check if a response is a validation error response.
 * A validation error response is an object with an "errors" property detailing the validation errors.
 *
 * Example:
 *   { "errors": { "email": "Email is required" } }
 */
function isValidationErrorResponse(
	data: unknown,
): data is ValidationErrorResponse {
	if (typeof data !== "object" || data === null) return false;
	if (!("errors" in data) || "error" in data) return false;
	if (typeof (data as { errors: Record<string, string> }).errors !== "object")
		return false;

	const errors = (data as { errors: Record<string, string> }).errors;
	return errors !== null && Object.keys(errors).length > 0;
}
