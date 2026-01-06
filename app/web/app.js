const itemsList = document.getElementById("items");
const itemTemplate = document.getElementById("item-template");
const emptyState = document.getElementById("empty-state");
const formStatus = document.getElementById("form-status");
const form = document.getElementById("item-form");
const refreshButton = document.getElementById("refresh");

async function loadItems() {
  toggleLoading(true);
  try {
    const response = await fetch("/api/items");
    if (!response.ok) throw new Error("Unable to load items");

    const data = await response.json();
    renderItems(data);
  } catch (err) {
    console.error(err);
    formStatus.textContent = "Failed to load items.";
  } finally {
    toggleLoading(false);
  }
}

function renderItems(items) {
  itemsList.innerHTML = "";
  emptyState.style.display = items.length ? "none" : "block";

  items.forEach((item) => {
    const node = itemTemplate.content.cloneNode(true);
    node.querySelector(".item-title").textContent = item.title;
    node.querySelector(".item-description").textContent = item.description || "No description";
    node.querySelector(".item-meta").textContent = new Date(item.createdAt).toLocaleString();

    const li = node.querySelector(".item");
    const deleteBtn = node.querySelector(".delete");
    deleteBtn.addEventListener("click", () => deleteItem(item.id, li));

    itemsList.appendChild(node);
  });
}

async function deleteItem(id, element) {
  const confirmed = confirm("Delete this item?");
  if (!confirmed) return;

  const response = await fetch(`/api/items/${id}`, { method: "DELETE" });
  if (!response.ok) {
    alert("Failed to delete item");
    return;
  }

  element.remove();
  if (!itemsList.children.length) {
    emptyState.style.display = "block";
  }
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  const formData = new FormData(form);
  const payload = {
    title: formData.get("title").trim(),
    description: formData.get("description").trim(),
  };

  if (!payload.title) {
    formStatus.textContent = "Title is required";
    return;
  }

  formStatus.textContent = "";
  form.querySelector("button[type='submit']").disabled = true;

  try {
    const response = await fetch("/api/items", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      const text = await response.text();
      throw new Error(text || "Unable to create item");
    }

    form.reset();
    await loadItems();
  } catch (err) {
    formStatus.textContent = err.message;
  } finally {
    form.querySelector("button[type='submit']").disabled = false;
  }
});

refreshButton.addEventListener("click", loadItems);

function toggleLoading(isLoading) {
  refreshButton.disabled = isLoading;
  refreshButton.textContent = isLoading ? "Loading..." : "Refresh";
}

loadItems();
