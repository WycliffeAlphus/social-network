import { getUser, getFollowers } from '../../utils/api';
import ProfileCard from '../../components/ProfileCard';
import FollowerCard from '../../components/FollowerCard';

export default function ProfilePage({ user, followers }) {
  return (
    <div>
      <ProfileCard user={user} />

      <h3 style={{ marginTop: '30px' }}>Followers</h3>
      <div>
        {followers.map(f => (
          <FollowerCard key={f.username} user={f} isFollowing={f.isFollowing} />
        ))}
      </div>
    </div>
  );
}

export async function getServerSideProps(context) {
  const { username } = context.params;

  const [user, followers] = await Promise.all([
    getUser(username),
    getFollowers(username)
  ]);

  return {
    props: { user, followers }
  };
}
