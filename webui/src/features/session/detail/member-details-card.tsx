import type { SessionMember } from "@/api/types.ts";
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { cn } from "@/lib/utils";

type MemberDetailsCardProps = {
	members: SessionMember[];
};

export function MemberDetailsCard({ members }: MemberDetailsCardProps) {
	return (
		<>
			<Card>
				<CardHeader>
					<CardTitle>Members</CardTitle>
				</CardHeader>
				<CardContent>
					{members.map((m, i) => (
						<div key={m.id}>
							<div className="flex items-center gap-4">
								<Avatar className="size-10">
									<AvatarFallback>{m.name[0]}</AvatarFallback>
								</Avatar>
								<div className="flex justify-between items-center w-full">
									<p className="py-4">{m.name}</p>
									<Balance n={-12} />
								</div>
							</div>
							{i < members.length - 1 && <Separator />}
						</div>
					))}
				</CardContent>
			</Card>
		</>
	);
}

function Balance({ n }: { n: number }) {
	return (
		<p
			className={cn(
				"text-muted-foreground",
				n > 0 && "text-green-500",
				n < 0 && "text-red-500",
			)}
		>
			{n}
		</p>
	);
}
