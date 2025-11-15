document.addEventListener("DOMContentLoaded", () => {
    const siteIdInput = document.getElementById("site-id");
    const dateInput = document.getElementById("date");
    const getStatsBtn = document.getElementById("get-stats-btn");
    const resultsContainer = document.getElementById("results-container");
    const errorContainer = document.getElementById("error-container");

    getStatsBtn.addEventListener("click", async () => {
        const siteId = siteIdInput.value;
        const date = dateInput.value;

        resultsContainer.innerHTML = "";
        errorContainer.innerHTML = "";

        if (!siteId || !date) {
            errorContainer.textContent = "Please enter both Site ID and Date.";
            return;
        }

        try {
            const response = await fetch(`/stats?site_id=${siteId}&date=${date}`);
            
            if (!response.ok) {
                const errData = await response.json();
                throw new Error(errData.error || `Error ${response.status}`);
            }

            const data = await response.json();

            let topPathsHTML = "<h4>Top Paths:</h4><ul>";
            if (data.top_paths.length > 0) {
                data.top_paths.forEach(path => {
                    topPathsHTML += `<li><strong>${path.path}</strong>: ${path.views} views</li>`;
                });
            } else {
                topPathsHTML += "<li>No page views recorded.</li>";
            }
            topPathsHTML += "</ul>";

            resultsContainer.innerHTML = `
                <h3>Stats for ${data.site_id} on ${data.date}</h3>
                <ul>
                    <li><strong>Total Views:</strong> ${data.total_views}</li>
                    <li><strong>Unique Users:</strong> ${data.unique_users}</li>
                </ul>
                ${topPathsHTML}
            `;

        } catch (err) {
            errorContainer.textContent = `Failed to get stats: ${err.message}`;
            console.error(err);
        }
    });
});