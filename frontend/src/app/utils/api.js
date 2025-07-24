const BASE_URL = "http://localhost:8080"; // your Go backend

export async function getUser(username) {
  const res = await fetch(`${BASE_URL}/api/users/${username}`);
  if (!res.ok) throw new Error("User not found");
  return res.json();
}

export async function getFollowers(username) {
  const res = await fetch(`${BASE_URL}/api/users/${username}/followers`);
  if (!res.ok) throw new Error("Failed to fetch followers");
  return res.json();
}
