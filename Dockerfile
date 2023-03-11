FROM jupyter/minimal-notebook:notebook-6.5.3

RUN rm -rf /home/jovyan/work && \
    pip install numpy pandas matplotlib scikit-learn
