import "../globals.css";

export default function ProfileCard({ user }) {
  return (
    <div className="profileCard">
      <img src={user.avatar || "/profile-placeholder.png"} alt="profile" />
      <h2>{user.name}</h2>
      <p>@{user.username}</p>
      <p>{user.bio}</p>
      <div className="stats">
        <span><strong>{user.followers}</strong> Followers</span>
        <span><strong>{user.following}</strong> Following</span>
      </div>
    </div>
  );
}
