async function Reaction(type, postid) {
  const res = await fetch(`http://localhost:8080/api/reaction?reaction_type=${type}&post_id=${postid}`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!res.ok) {
    throw new Error("Network response was not ok");
  }

  return res.json();
}

export default Reaction;