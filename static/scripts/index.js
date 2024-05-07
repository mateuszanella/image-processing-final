// The whole application relies on this 
// *-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-*
let file; // File object to be uploaded

async function refreshImage() {
    return new Promise((resolve, reject) => {
        var img = document.getElementById('image-display');
        fetch('./api/image?' + new Date().getTime()) // Prevent caching
            .then(response => response.blob())
            .then(blob => {
                img.src = URL.createObjectURL(blob);
                img.onload = resolve;
                img.onerror = reject;
            })
            .catch(reject);
    });
}

async function uploadImage(file) {
    let formData = new FormData();
    formData.append('image', file);

    return new Promise(async (resolve, reject) => {
        try {
            const response = await fetch('/api/image', {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            resolve();
        } catch (error) {
            reject(error);
        }
    });
}

// *-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-*

async function resetDropzone() {
    const label = document.querySelector('label[for="dropzone-file"]');
    label.classList.remove('border-solid', 'border-blue-500');
    label.classList.add('border-dashed');

    label.querySelector('p').classList.remove('hidden');

    label.querySelector('svg').classList.remove('text-blue-500');
    label.querySelector('svg').classList.add('text-gray-500', 'dark:text-gray-400');

    label.querySelector('p:last-child').classList.remove('hidden');

    label.querySelector('div').classList.remove('p-2');
    label.querySelector('div').classList.add('pt-5', 'pb-6');

    label.style.backgroundImage = '';
}

// *-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-*

// Image manipulation

// // *-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-*

// Filters
async function applyGrayscaleAndRefresh() {
    await fetch('/api/grayscale', { method: 'POST' });
    refreshImage();
}

async function applyBinaryAndRefresh() {
    var threshold = document.getElementById("threshold").value;
    await fetch('/api/binary', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ threshold: threshold }) 
    });
    refreshImage();
}

async function applyHistogramEqualizationAndRefresh() {
    await fetch('/api/histogram-equalization', { method: 'POST' });
    refreshImage();
}

async function applyNegativeAndRefresh() {
    await fetch('/api/negative', { method: 'POST' });
    refreshImage();
}

// Basic operations
async function applyAddAndRefresh() {
    var value = document.getElementById("add-value").value;
    await fetch('/api/add', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
            value: value 
        }) 
    });
    refreshImage();
}

async function applySubAndRefresh() {
    var value = document.getElementById("sub-value").value;
    await fetch('/api/subtract', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ value: value }) 
    });
    refreshImage();
}

async function applyMulAndRefresh() {
    var value = document.getElementById("mul-value").value;
    await fetch('/api/multiply', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ value: value }) 
    });
    refreshImage();
}

async function applyDivAndRefresh() {
    var value = document.getElementById("div-value").value;
    await fetch('/api/divide', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ value: value }) 
    });
    refreshImage();
}

// Logic operations
async function applyNOTAndRefresh() {
    await fetch('/api/not', { method: 'POST' });
    refreshImage();
}

// Spatial domain filters
async function applyMeanFilterAndRefresh() {
    var size = document.getElementById("mean-filter-size").value;
    await fetch('/api/mean-sdf', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ size: size }) 
    });
    refreshImage();
}