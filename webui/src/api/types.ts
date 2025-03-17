export type ValidationErrorResponse = {
	errors: Record<string, string>;
};

export function isValidationErrorResponse(
	data: any,
): data is ValidationErrorResponse {
	return (
		typeof data === "object" &&
		data !== null &&
		"errors" in data &&
		!("error" in data) &&
		typeof data.errors === "object" &&
		Object.values(data.errors).every((value) => typeof value === "string") &&
		Object.keys(data.errors).length > 0
	);
}
