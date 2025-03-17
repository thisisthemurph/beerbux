export type ValidationErrorResponse = {
	errors: Record<string, string>;
};

export function isValidationErrorResponse(
	data: unknown,
): data is ValidationErrorResponse {
	if (typeof data !== "object" || data === null) return false;
	if (!("errors" in data) || "error" in data) return false;
	if (typeof (data as { errors: Record<string, string> }).errors !== "object")
		return false;

	const errors = (data as { errors: Record<string, string> }).errors;
	return errors !== null && Object.keys(errors).length > 0;
}
