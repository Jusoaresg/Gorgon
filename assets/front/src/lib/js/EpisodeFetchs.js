import { PUBLIC_API_BASE_URL } from "$env/static/public";

export async function changeTrackingStatus(episodeID, newStatus) {
  try {
    const res = await fetch(
      `${PUBLIC_API_BASE_URL}/database/show/episode/status`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ episode_id: episodeID, tracking: newStatus }),
      },
    );

    if (!res.ok) {
      console.error("Failed to update episode status", res.status);
    }

    const data = await res.json();
    console.log("New Episode status", data);
  } catch (error) {
    console.error("Error trying to change episode status", error);
  }
}
