/* 
    The whole application relies on this 
    Yes, the image uploading and refreshing is the most unsafe thing ever seen by humanity
    But it works, and that's what matters, it's about the journey after all
*/
// *-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-*
let file; // File object to be uploaded
let fileType = 'jpeg'; // File type of the uploaded file

async function refreshImage() {
    return new Promise((resolve, reject) => {
        var img = document.getElementById('image-display');

        if (fileType === 'tiff' || fileType === 'tif') type = 'jpeg' 
        else type = fileType;

        fetch('./api/image?filetype=' + type +
            '&_=' + new Date().getTime()) // Prevent caching
            .then(response => response.blob())
            .then(blob => {
                img.src = URL.createObjectURL(blob);
                img.onload = resolve;
                img.onerror = reject;
            })
            .catch(reject);
    });
}

async function uploadImage(upload, filename = 'single') {
    let formData = new FormData();
    formData.append('image', upload);

    let apiEndpoint;
    switch(filename) {
        case 'single':
            apiEndpoint = '/api/image';
            break;
            case 'image1':
            case 'image2':
                apiEndpoint = '/api/image-combination';
                formData.append('filename', filename);
                break;
        default:
            throw new Error('Invalid operation type');
    }

    return new Promise(async (resolve, reject) => {
        try {
            const response = await fetch(apiEndpoint, {
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
    await fetch('/api/grayscale', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}

async function applyBinaryAndRefresh() {
    var threshold = document.getElementById("threshold").value;
    await fetch('/api/binary', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            threshold: threshold 
        }) 
    });
    refreshImage();
}

async function applyHistogramEqualizationAndRefresh() {
    await fetch('/api/histogram-equalization', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}

async function applyNegativeAndRefresh() {
    await fetch('/api/negative', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
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
            filename: 'uploaded.' + fileType,
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
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            value: value 
        }) 
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
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            value: value 
        }) 
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
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            value: value 
        }) 
    });
    refreshImage();
}

// Logic operations
async function applyNOTAndRefresh() {
    await fetch('/api/not', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
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
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            size: size 
        }) 
    });
    refreshImage();
}

async function applyMedianFilterAndRefresh() {
    var size = document.getElementById("median-filter-size").value;
    await fetch('/api/median-sdf', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            size: size 
        }) 
    });
    refreshImage();
}

async function applyGaussianFilterAndRefresh() {
    var size = document.getElementById("gaussian-filter-size").value;
    await fetch('/api/gaussian-sdf', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            size: size 
        }) 
    });
    refreshImage();
}

async function applyMinimumFilterAndRefresh() {
    var size = document.getElementById("minimum-filter-size").value;
    await fetch('/api/min-sdf', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            size: size 
        }) 
    });
    refreshImage();
}

async function applyMaximumFilterAndRefresh() {
    var size = document.getElementById("maximum-filter-size").value;
    await fetch('/api/max-sdf', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            size: size 
        }) 
    });
    refreshImage();
}

async function applyOrderFilterAndRefresh() {
    var position = document.getElementById("order-sdf-value").value;
    await fetch('/api/order-sdf', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            filename: 'uploaded.' + fileType,
            position: position
        }) 
    });
    refreshImage();
}

async function applyConservativeSmoothingAndRefresh() {
    await fetch('/api/conservative-smoothing-sdf', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}

// Morphological operations
async function applyDilationAndRefresh() {
    var size = document.getElementById("dilation-matrix-size").value;
    var kernelType = document.getElementById("dilation-kernel-type").value;
    await fetch('/api/dilation', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            size: size,
            kernelType: kernelType
        }) 
    });
    refreshImage();
}

async function applyErosionAndRefresh() {
    var size = document.getElementById("erosion-matrix-size").value;
    var kernelType = document.getElementById("erosion-kernel-type").value;
    await fetch('/api/erosion', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            size: size,
            kernelType: kernelType
        }) 
    });
    refreshImage();
}

async function applyOpeningAndRefresh() {
    var size = document.getElementById("opening-matrix-size").value;
    var kernelType = document.getElementById("opening-kernel-type").value;
    await fetch('/api/opening', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            size: size,
            kernelType: kernelType
        }) 
    });
    refreshImage();
}

async function applyClosingAndRefresh() {
    var size = document.getElementById("closing-matrix-size").value;
    var kernelType = document.getElementById("closing-kernel-type").value;
    await fetch('/api/closing', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
            filename: 'uploaded.' + fileType,
            size: size,
            kernelType: kernelType
        }) 
    });
    refreshImage();
}

async function applyContourAndRefresh() {
    await fetch('/api/contour', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}

//Edge detection
async function applyPrewittEdgeDetectionAndRefresh() {
    await fetch('/api/prewitt', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}

async function applySobelEdgeDetectionAndRefresh() {
    await fetch('/api/sobel', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}

async function applyLaplacianEdgeDetectionAndRefresh() {
    await fetch('/api/laplacian', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}

// *-**-* Bonus *-**-*

// Image Adjustments
async function flipLRAndRefresh() {
    await fetch('/api/flip-lr', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}

async function flipUDAndRefresh() {
    await fetch('/api/flip-ud', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}

async function rotate90AndRefresh() {
    await fetch('/api/rotate-90', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}

async function rotate270AndRefresh() {
    await fetch('/api/rotate-270', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ filename: 'uploaded.' + fileType }) 
    });
    refreshImage();
}