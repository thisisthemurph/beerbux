import NavigationDrawer from "@/components/navigation-drawer.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useBackNavigationStore } from "@/stores/back-navigation-store.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import { Link, useNavigate } from "react-router";

export function Navigation() {
	const user = useUserStore((state) => state.user);
	const backLink = useBackNavigationStore((state) => state.backLink);

	return (
		<nav className="flex justify-between items-center px-4 py-5 mb-4">
			{backLink ? (
				<BackLinkButton link={backLink} />
			) : (
				<Link to="/" className="text-xl tracking-wider">
					<span className="font-semibold">Beer</span>
					<span className="text-slate-600 dark:text-gray-400">bux</span>
				</Link>
			)}
			<div className="flex items-center gap-4">
				{!user && <LoginButton />}
				<NavigationDrawer />
			</div>
		</nav>
	);
}

function LoginButton() {
	return (
		<Button size="sm" variant="secondary" asChild className="rounded-full">
			<Link to="/login" className="text-sm tracking-wider">
				Login
			</Link>
		</Button>
	);
}

function BackLinkButton({ link }: { link: string }) {
	const navigate = useNavigate();

	return (
		<Button
			size="icon"
			variant="secondary"
			className="rounded-full"
			onClick={() => navigate(link)}
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				strokeWidth={1}
				stroke="currentColor"
				className="size-12"
			>
				<title>Back</title>
				{/* Circle */}
				<path
					className="stroke-none"
					d="M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
				/>
				{/* Arrow */}
				<path
					strokeLinecap="round"
					strokeLinejoin="round"
					d="M15 12H9m3-3-3 3 3 3"
				/>
			</svg>
		</Button>
	);
}
