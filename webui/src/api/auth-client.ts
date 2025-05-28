import { apiFetch } from "@/api/api-fetch.ts";
import type { User } from "./types/user.ts";

function useAuthClient() {
	const login = async (username: string, password: string): Promise<User> => {
		return apiFetch<User>("/auth/login", {
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

	const getCurrentUser = async (): Promise<User> => {
		return apiFetch<User>("/user");
	};

	return { login, signup, getCurrentUser };
}

export default useAuthClient;
