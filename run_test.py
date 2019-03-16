
import os
from shutil import copyfile
photo_id = "86016"

# Remove files in Optimal segmentation files
sol_folder = "./solutions/Optimal_Segmentation_Files"
for filename in os.listdir(sol_folder):
    file_path = os.path.join(sol_folder, filename)
    try:
        if os.path.isfile(file_path):
            os.unlink(file_path)
    except Exception as e:
        print(e)


# Copy files from solution folder
data_folder = "./data/" + photo_id
for filename in os.listdir("./data/" + photo_id):
    if "GT" in filename:
        file_path = os.path.join(data_folder, filename )
        dst = os.path.join(sol_folder, filename)
        copyfile(file_path, dst)
