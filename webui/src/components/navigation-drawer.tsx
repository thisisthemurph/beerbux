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
	DrawerTrigger,
} from "@/components/ui/drawer.tsx";
import { AlignRight, ChevronDown } from "lucide-react";
import { Link } from "react-router";

export default function NavigationDrawer() {
	return (
		<Drawer>
			<DrawerTrigger asChild className="px-0">
				<button type="button" title="open navigation">
					<AlignRight />
				</button>
			</DrawerTrigger>
			<DrawerContent>
				<DrawerHeader>
					<DrawerTitle className="sr-only">Navigation</DrawerTitle>
					<DrawerDescription className="sr-only">Navigation</DrawerDescription>
					<section className="flex flex-col gap-2 mt-6 text-center text-xl">
						<DrawerClose asChild>
							<Link to="/">Home</Link>
						</DrawerClose>
					</section>
				</DrawerHeader>
				<DrawerFooter>
					<section className="flex justify-between">
						<ThemeToggle />
						<DrawerClose asChild>
							<Button variant="outline" title="Close navigation">
								<ChevronDown />
							</Button>
						</DrawerClose>
					</section>
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	);
}
