import { useUser } from "@/hooks/useUser.tsx";

function HomePge() {
	const { user } = useUser();

	return (
		<div>
			<h1>Home Page</h1>
			{user ? <p>Welcome {user.username}</p> : <p>Welcome to Beerbux</p>}
		</div>
	);
}

export default HomePge;
