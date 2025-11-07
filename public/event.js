const eventId = new URLSearchParams(location.search).get("e");
const gallery = document.getElementById("gallery");
const fileInput = document.getElementById("file-input");
const uploadBtn = document.getElementById("upload-button");

let config = {
  allowGallery: true,
  allowAny: true,
  requireTakenToday: false,
};

let userId = localStorage.getItem("anonUserId");
if (!userId) {
  userId = crypto.randomUUID();
  localStorage.setItem("anonUserId", userId);
}

loadEventConfig();

// Load event rules from backend
async function loadEventConfig() {
  const res = await fetch(`/event-config?eventId=${eventId}`);
  if (res.ok) {
    config = await res.json();
    const titleEl = document.getElementById("event-title");
    if (titleEl && config.title) titleEl.textContent = config.title;
    applyConfig();
    loadImages(); // initial image load
    setInterval(loadImages, 5000); // refresh every 5s
  } else {
    alert("Event not found");
    uploadBtn.disabled = true;
  }
}

// Apply upload rules to input behavior
function applyConfig() {
  if (!config.allowGallery) {
    fileInput.setAttribute("capture", "environment");
  } else {
    fileInput.removeAttribute("capture");
  }
}

// Image upload logic
uploadBtn.onclick = async () => {
  const file = fileInput.files[0];
  if (!file) return alert("Choose a photo first.");

  if (!config.allowAny && config.requireTakenToday) {
    const isToday = await isTakenToday(file);
    if (!isToday) {
      alert("Only photos taken today are allowed.");
      return;
    }
  }

  const form = new FormData();
  form.append("image", file);
  form.append("eventId", eventId);
  form.append("userId", userId);

  const res = await fetch("/upload", { method: "POST", body: form });
  if (!res.ok) {
    alert("Upload failed");
    return;
  }

  fileInput.value = "";
  loadImages(); // reload after upload
};

// Check EXIF data for "DateTimeOriginal"
async function isTakenToday(file) {
  try {
    const exif = await exifr.parse(file, ['DateTimeOriginal']);
    if (!exif || !exif.DateTimeOriginal) return false;

    const takenDate = new Date(exif.DateTimeOriginal);
    const today = new Date();
    return (
      takenDate.getFullYear() === today.getFullYear() &&
      takenDate.getMonth() === today.getMonth() &&
      takenDate.getDate() === today.getDate()
    );
  } catch (e) {
    console.warn("EXIF read failed:", e);
    return false;
  }
}

// Render image gallery
async function loadImages() {
  gallery.innerHTML = "";

  const res = await fetch(`/images?eventId=${eventId}`);
  if (!res.ok) {
    console.error("Could not load images");
    return;
  }

  const images = await res.json();
  images.forEach((img) => {
    const wrapper = document.createElement("div");
    wrapper.className = "photo-card";

    const image = document.createElement("img");
    image.src = img.url;
    image.loading = "lazy";

    const label = document.createElement("small");
    if (img.userId === userId) {
      label.textContent = "You";
    } else if (img.userId) {
      label.textContent = `User ${img.userId.slice(0, 6)}`;
    } else {
      label.textContent = "Anonymous";
    }

    wrapper.appendChild(image);
    wrapper.appendChild(label);
    gallery.appendChild(wrapper);
  });
}
