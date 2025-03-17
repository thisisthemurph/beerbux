import {
	type ValidationErrorResponse,
	isValidationErrorResponse,
} from "@/api/types.ts";

type LoginResponse = {
	id: string;
	username: string;
};

function useAuthClient() {
	async function login(
		username: string,
		password: string,
	): Promise<LoginResponse | ValidationErrorResponse> {
		try {
			const response = await fetch("http://localhost:42069/auth/login", {
				method: "POST",
				credentials: "include",
				body: JSON.stringify({ username, password }),
				headers: {
					"Content-Type": "application/json",
				},
			});

			const data = await response.json();
			if (response.ok) {
				return data as LoginResponse;
			}

			if (isValidationErrorResponse(data)) {
				return data;
			}

			if ("error" in data) {
				throw Error(data.error ?? "An error occurred");
			}
			throw Error("An error occurred");
		} catch (error) {
			throw new Error(
				error instanceof Error ? error.message : "An error occurred",
			);
		}
	}

	async function signup(
		name: string,
		username: string,
		password: string,
		verificationPassword: string,
	): Promise<ValidationErrorResponse | undefined> {
		try {
			const response = await fetch("http://localhost:42069/auth/signup", {
				method: "POST",
				credentials: "include",
				body: JSON.stringify({
					name,
					username,
					password,
					verificationPassword,
				}),
				headers: {
					"Content-Type": "application/json",
				},
			});

			if (!response.ok) {
				const data = await response.json();
				if (isValidationErrorResponse(data)) {
					return data;
				}

				if ("error" in data) {
					throw Error(data.error ?? "An error occurred");
				}
				throw Error("An error occurred");
			}
		} catch (error) {
			throw new Error(
				error instanceof Error ? error.message : "An error occurred",
			);
		}
	}

	return { login, signup };
}

export default useAuthClient;
