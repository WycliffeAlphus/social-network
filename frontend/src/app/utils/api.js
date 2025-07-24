const BASE_URL = 'http://localhost:8080'; 

export async function getUser(username) {
  const res = await fetch(`${BASE_URL}/api/users/${username}`);
  return res.json();
}

export async function getFollowers(username) {
  const res = await fetch(`${BASE_URL}/api/users/${username}/followers`);
  return res.json();
}
