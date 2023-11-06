# iforest_go_MNIST

## Summary
This project uses MNIST dataset, prepare the dataset and fit processed data to isolation forest model. The model create output of anomaly score. At last, create a correlation calculation that comparing with Python and R separately. 
To process the MNIST dataset, the project use GoMNIST package to process teh image data which contains 60000 training image and 10000 for model testing. In this study, only the training data has been participated. 

## Package selection
The package I found the most usful and have the best approach to isolation forest is the [e-XpertSolutions/go-iforest](https://github.com/e-XpertSolutions/go-iforest/tree/master#readme) package. Besides the go-iforest, there are other several considerable package such as codegaudi/go-iforest, malaschitz/randomForest/isolationForest. The package I chose provide great example and straightforward in anomaly score and prediction. This package can satisfy most of the testing as it created just for isolation forest. No customized data type and easy to use. 

## In this project
*How easy was it to clone the package repository, implement the code, and set up the tests?  Were the tests successful?  To what extent were Go results similar to results from R and/or Python? How difficult will it be for the firm to implement Go-based isolation forests as part of its data processing pipeline?*

The implementation is difficult even though the image processing and model fitting was easier. The difficulty for is mostly coming from dealing with data structure as initially I was using dataframe but find this could lead to more complexity of code. The test result is stored in csv that range in 0-0.2. Comparing to manipulate data, finding ideal package and model fitting, prediciton, Python and R have very developed library. As a firm intend to using Go for data processing, model creating, the task can be difficult by setting up and processing speed. In this study, I also feel the performance of Go is not outstanding and even slower when increase the tree number let alone creating a hyperparameter tuning study. 

## Correlation
At the end of the program, the correlation score was calculated against Python and R separately and indicates an inverse relationship which need further attention as this is very big difference comparing to Python and R. For further study, it is necessary to try out multiple isolation forest package to compare the difference which absent from this project. 

Python vs Go-iforest: The correlation score is -0.94

R vs Go-iforest: The correlation score is -0.68

