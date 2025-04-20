import {
	Drawer,
	DrawerContent,
	DrawerDescription,
	DrawerHeader,
	DrawerTitle,
} from "@/components/ui/drawer.tsx";
import { AddMemberForm } from "./add-member-form";

interface DrawerToggleProps {
	open: boolean;
	onOpenChange: (open: boolean) => void;
}

interface AddMemberDrawerProps extends DrawerToggleProps {
	onMemberAdd: (username: string) => void;
}

export function AddMemberDrawer({
	onMemberAdd,
	...drawerToggleProps
}: AddMemberDrawerProps) {
	return (
		<Drawer {...drawerToggleProps}>
			<DrawerContent>
				<DrawerHeader>
					<section className="flex justify-between items-center">
						<DrawerTitle>Add a new member</DrawerTitle>
					</section>
					<DrawerDescription>
						Add a new member to the session by their username.
					</DrawerDescription>
				</DrawerHeader>
				<section className="p-4">
					<AddMemberForm onAdd={onMemberAdd} />
				</section>
			</DrawerContent>
		</Drawer>
	);
}
