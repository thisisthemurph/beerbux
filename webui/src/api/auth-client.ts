import { apiFetch } from "@/api/api-fetch.ts";
import type { UserAuthDetails } from "./types/user.ts";

function useAuthClient() {
	const login = async (username: string, password: string): Promise<UserAuthDetails> => {
		return apiFetch<UserAuthDetails>("/auth/login", {
			method: "POST",
			body: JSON.stringify({ username, password }),
		});
	};

	const signup = async (
		name: string,
		username: string,
		email: string,
		password: string,
		verificationPassword: string,
	): Promise<void> => {
		return apiFetch<void>("/auth/signup", {
			method: "POST",
			body: JSON.stringify({ name, username, email, password, verificationPassword }),
		});
	};

	const getCurrentUser = async (): Promise<UserAuthDetails> => {
		return apiFetch<UserAuthDetails>("/user");
	};

	const initializePasswordReset = async (newPassword: string): Promise<void> => {
		return apiFetch<void>("/auth/password/initialize-update", {
			method: "POST",
			body: JSON.stringify({ newPassword }),
		});
	};

	const resetPassword = async (otp: string): Promise<void> => {
		return apiFetch<void>("/auth/password", {
			method: "PUT",
			body: JSON.stringify({ otp }),
		});
	};

	const initializeEmailUpdate = async (newEmail: string): Promise<void> => {
		return apiFetch<void>("/auth/email/initialize-update", {
			method: "POST",
			body: JSON.stringify({ newEmail }),
		});
	};

	const updateEmail = async (otp: string): Promise<void> => {
		return apiFetch<void>("/auth/email", {
			method: "PUT",
			body: JSON.stringify({ otp }),
		});
	};

	return {
		login,
		signup,
		getCurrentUser,
		initializePasswordReset,
		resetPassword,
		initializeEmailUpdate,
		updateEmail,
	};
}

export default useAuthClient;
