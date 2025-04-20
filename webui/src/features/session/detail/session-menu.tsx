import { Button } from "@/components/ui/button.tsx";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useConfirmationDialog } from "@/hooks/use-confirmation-dialog";
import {
	ArrowLeftFromLine,
	CircleX,
	EllipsisVertical,
	PauseCircle,
} from "lucide-react";

export function SessionMenu({
	showAdminActions,
}: { showAdminActions: boolean }) {
	const [openConfirmationDialog, ConfirmationDialog] = useConfirmationDialog();

	const handleCloseSession = () => {
		openConfirmationDialog({
			title: "Are you sure you want to close this session?",
			description:
				"Are you sure you want to close this session? Once closed, members will be able to see the session, but will never be able to interact with it again. Once closed, the session cannot be reopened.",
			confirmText: "Close session",
			cancelText: "Cancel",
			onConfirm: () => {
				console.log("close session not implemented");
			},
		});
	};

	const handlePauseSession = () => {
		openConfirmationDialog({
			title: "Are you sure you want to pause this session?",
			description:
				"Pausing the session will prevent anyone from tracking any new rounds being purchased until the session is unpaused. " +
				"Any admin user can unpause the session at any time.",
			confirmText: "Pause session",
			cancelText: "Cancel",
			onConfirm: () => {
				console.log("pause session not implemented");
			},
		});
	};

	const handleLeaveSession = () => {
		openConfirmationDialog({
			title: "Are you sure you would like to leave?",
			description:
				"After leaving the session, if you would like to rejoin, you will need to be invited again.",
			confirmText: "Leave session",
			cancelText: "Stay",
			onConfirm: () => {
				console.log("leave session not implemented");
			},
		});
	};

	return (
		<>
			<ConfirmationDialog />
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button size="icon" variant="secondary" className="rounded-full">
						<EllipsisVertical />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent className="min-w-64 mx-4">
					{showAdminActions && (
						<>
							<DropdownMenuGroup>
								<DropdownMenuLabel className="text-lg font-semibold">
									Admin actions
								</DropdownMenuLabel>
								<DropdownMenuSeparator />
								<DropdownMenuItem
									className="group text-lg gap-4 cursor-pointer"
									onClick={handleCloseSession}
								>
									<CircleX className="size-6 group-hover:text-destructive/75 transition-colors" />
									<span>Close session</span>
								</DropdownMenuItem>
								<DropdownMenuItem
									className="group text-lg gap-4 cursor-pointer"
									onClick={handlePauseSession}
								>
									<PauseCircle className="size-6 group-hover:text-destructive/75 transition-colors" />
									<span>Pause session</span>
								</DropdownMenuItem>
							</DropdownMenuGroup>
							<DropdownMenuSeparator />
						</>
					)}
					<DropdownMenuGroup>
						<DropdownMenuItem
							className="group text-lg gap-4 cursor-pointer"
							onClick={handleLeaveSession}
						>
							<ArrowLeftFromLine className="size-6 text-muted-foreground group-hover:text-primary/75 transition-colors" />
							<span>Leave session</span>
						</DropdownMenuItem>
					</DropdownMenuGroup>
				</DropdownMenuContent>
			</DropdownMenu>
		</>
	);
}

SessionMenu.Skeleton = function SessionMenuSkeleton() {
	return (
		<Button
			size="icon"
			variant="secondary"
			className="rounded-full animate-pulse"
		>
			<EllipsisVertical />
		</Button>
	);
};
