import {
	Drawer,
	DrawerContent,
	DrawerDescription,
	DrawerHeader,
	DrawerTitle,
} from "@/components/ui/drawer.tsx";
import {
	CreateSessionForm,
	type CreateSessionFormValues,
} from "@/features/dashboard/create-session/create-session-form.tsx";

interface DrawerToggleProps {
	open: boolean;
	onOpenChange: (open: boolean) => void;
}

interface CreateSessionDrawerProps extends DrawerToggleProps {
	onCreate: (values: CreateSessionFormValues) => void;
}

export function CreateSessionDrawer({ onCreate, ...drawerToggleProps }: CreateSessionDrawerProps) {
	return (
		<Drawer {...drawerToggleProps}>
			<DrawerContent>
				<DrawerHeader>
					<section className="flex justify-between items-center">
						<DrawerTitle>Start a new session</DrawerTitle>
					</section>
					<DrawerDescription>
						Create a new session for you and your fiends. You can add your friends once the session has been
						created.
					</DrawerDescription>
				</DrawerHeader>
				<section className="p-4">
					<CreateSessionForm onCreate={onCreate} />
				</section>
			</DrawerContent>
		</Drawer>
	);
}
