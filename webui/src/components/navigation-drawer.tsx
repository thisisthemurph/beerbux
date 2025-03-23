import useUserClient from "@/api/user-client.ts";
import { ThemeToggle } from "@/components/theme-toggle.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
	Drawer,
	DrawerClose,
	DrawerContent,
	DrawerDescription,
	DrawerFooter,
	DrawerHeader,
	DrawerTitle,
} from "@/components/ui/drawer.tsx";
import { cn } from "@/lib/utils.ts";
import {
	type NavigationView,
	useNavigationStore,
} from "@/stores/navigation-store.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import {
	AlignRight as BurgerMenuIcon,
	ChevronDown,
	ChevronLeft,
	LogOut,
	User2,
} from "lucide-react";
import type * as React from "react";
import { Link } from "react-router";
import { toast } from "sonner";

export default function NavigationDrawer() {
	const { logout } = useUserClient();
	const localLogout = useUserStore((state) => state.logout);
	const { view, setView, isOpen, open, close } = useNavigationStore();
	const user = useUserStore((state) => state.user);

	function handleLogout() {
		logout()
			.then(localLogout)
			.then(close)
			.then(() => toast.success("Logged out successfully"))
			.catch(() => toast.error("Failed to logout"));
	}

	return (
		<Drawer open={isOpen} onClose={close}>
			{user && <UserButton onClick={() => open("user")} />}
			<button
				type="button"
				title="Open navigation drawer"
				onClick={() => open()}
			>
				<BurgerMenuIcon />
			</button>
			<DrawerContent>
				<DrawerHeader>
					<DrawerTitle
						className={cn(
							"sr-only text-center font-mono tracking-wider text-2xl",
							view !== "primary" && "not-sr-only",
						)}
					>
						{user?.username}
					</DrawerTitle>
					<DrawerDescription className="sr-only">Navigation</DrawerDescription>
					<section className="flex flex-col gap-2 my-6 text-center text-xl">
						<NavigationContent view={view} />
					</section>
				</DrawerHeader>
				<DrawerFooter>
					<section className="flex items-center justify-between">
						<div className="space-x-2">
							{/* Show the user button when the user is logged in */}
							{view === "primary" && user && (
								<Button variant="secondary" onClick={() => setView("user")}>
									<User2 />
								</Button>
							)}

							{/* Sow the theme toggle when the user is logged out */}
							{!user && <ThemeToggle />}

							{/* Show the back button and logout button when the user is viewing their profile */}
							{view === "user" && user && (
								<>
									<Button
										variant="secondary"
										title="Back to main navigation"
										onClick={() => setView("primary")}
									>
										<ChevronLeft />
									</Button>
									<Button variant="secondary" onClick={handleLogout}>
										<LogOut />
										<span>Logout</span>
									</Button>
								</>
							)}
						</div>

						{/* Show the close navigation button at all times */}
						<div className="space-x-2">
							{/* Show the theme toggle when the user is logged in and in user view */}
							{view === "user" && user && <ThemeToggle />}
							<DrawerClose asChild>
								<Button
									variant="outline"
									title="Close navigation"
									onClick={() => setView("primary")}
								>
									<ChevronDown />
								</Button>
							</DrawerClose>
						</div>
					</section>
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	);
}

function NavigationContent({ view }: { view: NavigationView }) {
	switch (view) {
		case "primary":
			return <PrimaryNavigationContent />;
		case "user":
			return <UserNavigationContent />;
	}
}

function PrimaryNavigationContent() {
	return (
		<DrawerClose asChild>
			<Link to="/">Home</Link>
		</DrawerClose>
	);
}

function UserNavigationContent() {
	return (
		<>
			<Link to="/profile">Profile</Link>
			<Link to="/settings">Settings</Link>
		</>
	);
}

function UserButton(props: React.ComponentProps<"button">) {
	return (
		<Button size="icon" variant="secondary" className="rounded-full" {...props}>
			<User2 />
		</Button>
	);
}
