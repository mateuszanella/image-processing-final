document.getElementById('processBtn').addEventListener('click', handleImageInput);

function handleImageInput() {
    var fileInput = document.getElementById('imageInputForm');
    var inputImg = document.getElementById('inputImage');
    var outputImg = document.getElementById('outputImage');
    
    if (fileInput.files && fileInput.files[0]) {
        var reader = new FileReader();
        
        reader.onload = function(e) {
            inputImg.src = e.target.result;
            outputImg.src = e.target.result;
        }
        
        reader.readAsDataURL(fileInput.files[0]);
    }
}