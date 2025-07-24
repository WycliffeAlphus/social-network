import ProfileCard from "@/app/components/ProfileCard";
import FollowerCard from "@/app/components/FollowerCard";
import { getUser, getFollowers } from "@/app/utils/api";

export default async function ProfilePage({ params }) {
  const username = params.username;

  // ✅ Fetch real data from your Go backend
  const [user, followers] = await Promise.all([
    getUser(username),
    getFollowers(username),
  ]);

  return (
    <div style={{ maxWidth: "600px", margin: "0 auto", padding: "20px" }}>
      <ProfileCard user={user} />

      <h3 style={{ marginTop: "30px" }}>Followers</h3>
      <div>
        {followers.map((f) => (
          <FollowerCard key={f.username} user={f} isFollowing={f.isFollowing} />
        ))}
      </div>
    </div>
  );
}
