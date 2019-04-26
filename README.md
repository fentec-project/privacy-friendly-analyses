# Privacy-Friendly Analyses Services

This repository demonstrates a centralized Functional Encryption setup (an alternative is a decentralized version where no central authority is needed for key generation).

![Privacy-Friendly Analyses components](https://raw.githubusercontent.com/fentec-project/private-analyses/master/img/CDV.png)

Demonstrator comprises three components:

 * Key Server: central authority component which generates keys.

 * Privacy-Friendly Analyses as a Service: the user sends encrypted data to this component and obtains the result of analysis computed using only encrypted data. Computation is enabled by Functional Encryption keys which are obtained from the Key Server.

 * Client: a component which obtains the public key from the Key Server, encrypts
user’s data with the public key and sends it to the Privacy-Friendly Analyses as a
Service component.

Various analysis services will be added in the future. Currently, Privacy-Friendly Analyses component contains a service to compute the 30-year risk of a general cardiovascular disease (CDV) based on the algorithm from [1], developed using the Framingham heart study [2].

The Framingham heart study followed roughly 5,000 patients from Framingham, Massachusettes, for many decades starting in 1948. Later, other patients were included. The risk models are algorithms used to assess the risk of specific atherosclerotic CDV events (coronary heart disease, cerebrovascular disease, peripheral vascular disease, heart failure). Algorithms most often estimate the 10-year or 30-year CDV risk of an individual.

The input parameters for algorithms are sex, age, total and high-density
lipoprotein cholesterol, systolic blood pressure, treatment for hypertension, smoking, and diabetes status. The demonstrator shows how the risk score can be computed using only the encrypted values of the input parameters. 

The user specifies the parameters in the Client program, and these are encrypted and sent to the Privacy-Friendly Analyses as a Service component. Privacy-Friendly Analyses component computes the 30-year risk and returns it to the user.

```
x := data.NewVector([]*big.Int{isMaleInt, ageInt, systolicBPInt, totalChInt, hdlChInt, smokerInt, treatedBPInt, diabeticInt})
```

Framingham risk score algorithms are based on Cox proportional hazards
model [3]. Part of it is multiplication of the input parameters by regression
factors which are real numbers. In 30-year algorithm the vector x is multiplied
by two vectors (scalar or inner-product):

```
y1 = (0.34362, 2.63588, 1.8803, 1.12673, -0.90941, 0.59397, 0.5232,
0.68602)
y2 = (0.48123, 3.39222, 1.39862, -0.00439, 0.16081, 0.99858,
0.19035, 0.49756)
```

Regression factors need to be converted into integers because cryptographic
schemes operate with integers. Thus, we multiply factors by the power of 10 to obtain whole numbers. A factor of 100 000 is used in our case. Consequently, we multiply the input parameters by the same factor. For example, boolean parameters which are presented as 1 (true) or 0 (false) thus become 100 000 or 0.

Client encrypts vector x using public key obtained from the Key Server:

```
ciphertext, err := paillier.Encrypt(x, masterPubKey)
```

Client then sends ciphertext to the Privacy-Friendly Analyses component. Privacy-Friendly Analyses beforehand obtained two functional encryption keys from the Key
Server: key to compute the inner-product of x and y1, and key to compute the
inner-product of x and y2. Now it can compute the inner-products:

```
xy1, err := paillier.Decrypt(ciphertext, key1, y1)
xy2, err := paillier.Decrypt(ciphertext, key2, y2)
```

To obtain the actual values of inner-products both values need to be divided
by 100 000 * 100 000. The algorithm then uses these two values to compute the risk. The risk value is returned to the Client. 

Note that the Privacy-Friendly Analyses component knows the risk, but does not know the user's data (sex, age, ... ).




[1] Pencina, M.J., D’Agostino Sr, R.B., Larson, M.G., Massaro, J.M., Vasan, R.S.:
Predicting the thirty-year risk of cardiovascular disease: the framingham heart study. Circulation 119(24), 3078 (2009)

[2] Dagostino, R.B., Vasan, R.S., Pencina, M.J., Wolf, P.A., Cobain, M., Massaro, J.M., Kannel, W.B.: General cardiovascular risk profile for use in primary care. Circulation 117(6), 743–753 (2008)

[3] Cox, D.R.: Regression models and life-tables. Journal of the Royal Statistical Society: Series B (Methodological) 34(2), 187–202 (1972)

