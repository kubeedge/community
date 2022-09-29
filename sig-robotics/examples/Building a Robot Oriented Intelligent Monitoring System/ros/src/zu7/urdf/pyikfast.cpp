#define IKFAST_HAS_LIBRARY

#include <Python.h>
#include <vector>
#include "ikfast.h"

using namespace ikfast;

static PyObject *inverse(PyObject *self, PyObject *args);
static PyObject *forward(PyObject *self, PyObject *args);
bool ComputeIk(const IkReal *eetrans, const IkReal *eerot, const IkReal *pfree, IkSolutionListBase<IkReal> &solutions);
void ComputeFk(const IkReal *j, IkReal *eetrans, IkReal *eerot);

static PyMethodDef PyIkFastMethods[] = {
    {"inverse", inverse, METH_VARARGS, "Calculate inverse kinematics"},
    {"forward", forward, METH_VARARGS, "Calculate forwards kinematics"},
    {NULL, NULL, 0, NULL}};

static struct PyModuleDef pyikfastmodule = {
    PyModuleDef_HEAD_INIT,
    "pyikfast",
    "ikfast wrapper",
    -1,
    PyIkFastMethods};

PyMODINIT_FUNC PyInit_pyikfast(void)
{
  PyObject *module = PyModule_Create(&pyikfastmodule);
  return module;
}

PyObject *inverse(PyObject *self, PyObject *args)
{
  PyObject *argTranslation;
  PyObject *argRotation;
  IkReal rotation[9];
  IkReal translation[3];

  // Parse arguments
  if (!PyArg_ParseTuple(args, "OO", &argTranslation, &argRotation))
  {
    return NULL;
  }
  for (int i = 0; i < 3; i++)
  {
    translation[i] = PyFloat_AsDouble(PyList_GetItem(argTranslation, i));
  }
  for (int i = 0; i < 9; i++)
  {
    rotation[i] = PyFloat_AsDouble(PyList_GetItem(argRotation, i));
  }

  // Compute inverse kinematics
  IkSolutionList<IkReal> solutions;
  ComputeIk(translation, rotation, NULL, solutions);

  // Return the solution
  PyObject *pySolutionCollection = PyList_New((int)solutions.GetNumSolutions());
  std::vector<IkReal> solvalues(GetNumJoints());
  for (int i = 0; i < solutions.GetNumSolutions(); i++)
  {
    const IkSolutionBase<IkReal> &sol = solutions.GetSolution(i);
    std::vector<IkReal> vsolfree(sol.GetFree().size());
    sol.GetSolution(&solvalues[0], vsolfree.size() > 0 ? &vsolfree[0] : NULL);

    PyObject *pySolution = PyList_New(solvalues.size());
    for (int j = 0; j < solvalues.size(); j++)
    {
      PyList_SetItem(pySolution, j, PyFloat_FromDouble(solvalues[j]));
    }
    PyList_SetItem(pySolutionCollection, i, pySolution);
  }

  return pySolutionCollection;
}

PyObject *forward(PyObject *self, PyObject *args)
{
  PyObject *argPositions;
  IkReal positions[20];
  IkReal rotation[9];
  IkReal translation[3];

  // Parse arguments
  if (!PyArg_ParseTuple(args, "O", &argPositions))
  {
    return NULL;
  }
  for (int i = 0; i < PyList_Size(argPositions); i++)
  {
    positions[i] = PyFloat_AsDouble(PyList_GetItem(argPositions, i));
  }

  // Compute forward kinematics
  ComputeFk(positions, translation, rotation);

  // Return the solution
  PyObject *pyResult = PyList_New(2);
  PyObject *pyTranslation = PyList_New(3);
  PyObject *pyRotation = PyList_New(9);
  PyList_SetItem(pyResult, 0, pyTranslation);
  PyList_SetItem(pyResult, 1, pyRotation);
  for (int i = 0; i < 3; i++)
  {
    PyList_SetItem(pyTranslation, i, PyFloat_FromDouble(translation[i]));
  }
  for (int i = 0; i < 9; i++)
  {
    PyList_SetItem(pyRotation, i, PyFloat_FromDouble(rotation[i]));
  }
  return pyResult;
}