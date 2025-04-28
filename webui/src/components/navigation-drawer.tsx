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
import { useNavigationStore } from "@/stores/navigation-store.ts";
import { useUserStore } from "@/stores/user-store.tsx";
import {
	Bolt,
	AlignRight as BurgerMenuIcon,
	ChevronDown,
	Home,
	LogIn,
	LogOut,
	User2,
	UserRoundPen,
} from "lucide-react";
import type * as React from "react";
import { Link, type LinkProps } from "react-router";
import { toast } from "sonner";

export default function NavigationDrawer() {
	const nav = useNavigationStore();
	const { logout } = useUserClient();
	const { user, isLoggedIn, logout: localLogout } = useUserStore();

	function handleLogout() {
		logout()
			.then(localLogout)
			.then(nav.close)
			.then(() => toast.success("Logged out successfully"))
			.catch(() => toast.error("Failed to logout"));
	}

	return (
		<Drawer open={nav.isOpen} onClose={nav.close} onOpenChange={nav.toggle}>
			{user ? <LoggedInNavButton onClick={nav.open} /> : <LoggedOutNavButton onClick={nav.open} />}
			<DrawerContent>
				<DrawerHeader>
					<DrawerTitle className="text-center font-mono tracking-wider text-2xl">
						{isLoggedIn ? user.username : "Beerbux"}
					</DrawerTitle>
					<DrawerDescription className="sr-only">Navigation</DrawerDescription>
				</DrawerHeader>

				<nav className="grid grid-cols-2 gap-2 mx-2 mb-8">
					{isLoggedIn ? <LoggedInNavigationMenu username={user?.username} /> : <LoggedOutNavigationMenu />}
				</nav>

				<DrawerFooter className="flex flex-row items-center justify-between">
					{isLoggedIn ? <LogoutButton handleLogout={handleLogout} /> : <LoginButton />}

					<div className="space-x-2">
						<ThemeToggle />
						<DrawerClose asChild>
							<Button variant="outline" title="Close navigation">
								<ChevronDown />
							</Button>
						</DrawerClose>
					</div>
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	);
}

function LoggedInNavButton(props: React.ComponentProps<"button">) {
	return (
		<Button size="icon" variant="secondary" className="rounded-full" {...props}>
			<User2 />
		</Button>
	);
}

function LoggedOutNavButton(props: React.ComponentProps<"button">) {
	return (
		<button type="button" {...props}>
			<BurgerMenuIcon />
		</button>
	);
}

function LoggedOutNavigationMenu() {
	return (
		<>
			<NavCloseLink to="/">Home</NavCloseLink>
			<NavCloseLink to="/about">About</NavCloseLink>
		</>
	);
}

function LoggedInNavigationMenu({ username }: { username: string }) {
	return (
		<>
			<NavCloseLink to="/">
				<Home />
				<span>{username}</span>
			</NavCloseLink>
			<NavCloseLink to="/profile">
				<UserRoundPen />
				<span>Profile</span>
			</NavCloseLink>
			<NavCloseLink to="/settings">
				<Bolt />
				<span>Settings</span>
			</NavCloseLink>
		</>
	);
}

function LogoutButton({ handleLogout }: { handleLogout: () => void }) {
	return (
		<Button variant="secondary" onClick={handleLogout}>
			<LogOut />
			<span>Logout</span>
		</Button>
	);
}

function LoginButton() {
	return (
		<DrawerClose asChild>
			<Button variant="secondary" asChild>
				<Link to="/login">
					<LogIn />
					<span>Login</span>
				</Link>
			</Button>
		</DrawerClose>
	);
}

function NavCloseLink({ children, className, ...props }: LinkProps) {
	return (
		<DrawerClose asChild>
			<Link
				{...props}
				className={cn(
					className,
					"flex justify-center items-center gap-2 p-6 font-semibold tracking-wide border border-primary/50 dark:border-primary/30 rounded hover:bg-primary/50 dark:hover:bg-primary/30 transition-colors",
				)}
			>
				{children}
			</Link>
		</DrawerClose>
	);
}
