import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button.tsx";
import { type JSX, useCallback, useState } from "react";

interface ConfirmationDialogOptions {
	title: string;
	description?: string;
	confirmText?: string;
	cancelText?: string;
	onConfirm: () => void;
}

export function useConfirmationDialog(): [
	(options: ConfirmationDialogOptions) => void,
	() => JSX.Element | null,
] {
	const [open, setOpen] = useState(false);
	const [options, setOptions] = useState<ConfirmationDialogOptions | null>(
		null,
	);

	const openDialog = useCallback((opts: ConfirmationDialogOptions) => {
		setOptions(opts);
		setOpen(true);
	}, []);

	const handleConfirm = useCallback(() => {
		if (options?.onConfirm) {
			options.onConfirm();
		}
		setOpen(false);
	}, [options]);

	const handleCancel = useCallback(() => {
		setOpen(false);
	}, []);

	const ConfirmationDialog = () => {
		if (!options) return null;
		return (
			<AlertDialog open={open} onOpenChange={setOpen}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle className="text-left text-2xl">
							{options.title}
						</AlertDialogTitle>
						{options.description && (
							<AlertDialogDescription className="text-left text-lg">
								{options.description}
							</AlertDialogDescription>
						)}
					</AlertDialogHeader>
					<AlertDialogFooter>
						<section className="flex gap-2 justify-end tracking-wider">
							<AlertDialogCancel onClick={handleCancel} asChild>
								<Button variant="secondary" size="lg">
									{options.cancelText ?? "Cancel"}
								</Button>
							</AlertDialogCancel>
							<AlertDialogAction onClick={handleConfirm} asChild>
								<Button variant="secondary" size="lg">
									{options.confirmText ?? "Confirm"}
								</Button>
							</AlertDialogAction>
						</section>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>
		);
	};

	return [openDialog, ConfirmationDialog];
}
