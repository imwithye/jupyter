FROM jupyter/minimal-notebook:notebook-6.5.3

RUN rm -rf /home/jovyan/work && \
    # Basic Python packages
    pip install ipywidgets numpy pandas matplotlib scikit-learn && \
    # PyTorch
    pip install torch torchvision torchaudio
