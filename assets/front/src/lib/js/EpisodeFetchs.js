export async function changeTrackingStatus(episodeID, newStatus) {
  try {
    const res = await fetch("api/v1/database/show/episode/status", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ episode_id: episodeID, tracking: newStatus }),
    });

    if (!res.ok) {
      console.error("Failed to update episode status", res.status);
    }

    const data = await res.json();
    console.log("New Episode status", data);
  } catch (error) {
    console.error("Error trying to change episode status", error);
  }
}
