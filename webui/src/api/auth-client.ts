import { apiFetch } from "@/api/api-fetch.ts";

type LoginResponse = {
	id: string;
	username: string;
};

function useAuthClient() {
	const login = async (
		username: string,
		password: string,
	): Promise<LoginResponse> => {
		return apiFetch<LoginResponse>("/auth/login", {
			method: "POST",
			body: JSON.stringify({ username, password }),
		});
	};

	const signup = async (
		name: string,
		username: string,
		password: string,
		verificationPassword: string,
	): Promise<void> => {
		return apiFetch<void>("/auth/signup", {
			method: "POST",
			body: JSON.stringify({ name, username, password, verificationPassword }),
		});
	};

	return { login, signup };
}

export default useAuthClient;
