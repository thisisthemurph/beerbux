import type { SessionMember } from "@/api/types.ts";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@/components/ui/card.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { UserAvatar } from "@/components/user-avatar.tsx";
import type { AvatarData } from "@/hooks/user-avatar-data.ts";
import { cn } from "@/lib/utils";

type MemberDetailsCardProps = {
	members: SessionMember[];
	avatarData: Record<string, AvatarData>;
};

export function MemberDetailsCard({
	members,
	avatarData,
}: MemberDetailsCardProps) {
	return (
		<Card>
			<CardHeader>
				<CardTitle>Members</CardTitle>
			</CardHeader>
			<CardContent>
				{members.map((m, i) => (
					<div key={m.id}>
						<div className="flex items-center gap-4">
							<UserAvatar data={avatarData[m.username]} tooltip={m.name} />
							<div className="flex justify-between items-center w-full">
								<p className="py-6">{m.username}</p>
								<Balance n={-12} />
							</div>
						</div>
						{i < members.length - 1 && <Separator />}
					</div>
				))}
			</CardContent>
		</Card>
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
			${n < 0 ? n * -1 : n}
		</p>
	);
}
