import "../globals.css";

export default function FollowerCard({ user, isFollowing }) {
  return (
    <div className="followerCard">
      <img src={user.avatar || "/profile-placeholder.png"} alt="avatar" />
      <div>
        <h4>{user.name}</h4>
        <p>@{user.username}</p>
      </div>
      <button className={isFollowing ? "following" : "follow"}>
        {isFollowing ? "Following" : "Follow"}
      </button>
    </div>
  );
}
