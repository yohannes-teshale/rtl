const fs = require('fs');
const path = require('path');
const parser = require('@babel/parser');
const traverse = require('@babel/traverse').default;

function resolveImportPath(importPath, importerPath, baseDir, testDir) {
    if (importPath.startsWith('.')) {
        const resolvedPath = path.resolve(path.dirname(importerPath), importPath);
        let filePath = resolvedPath;

        if (!fs.existsSync(filePath)) {
            filePath = fs.existsSync(`${resolvedPath}.ts`) ? `${resolvedPath}.ts` :
                fs.existsSync(`${resolvedPath}.js`) ? `${resolvedPath}.js` :
                    fs.existsSync(`${resolvedPath}/index.ts`) ? `${resolvedPath}/index.ts` :
                        fs.existsSync(`${resolvedPath}/index.js`) ? `${resolvedPath}/index.js` : null;
        }
        return filePath.startsWith(testDir) ? filePath : null;
    }
    return null;
}

function parseFile(filePath) {
    console.log("parsing file",filePath);
    try{
        const content = fs.readFileSync(filePath, 'utf8');
        return parser.parse(content, {
            sourceType: 'module',
            plugins: ['jsx', 'classProperties','typescript'],
        });
    }catch(e){
        console.error("Error parsing file",filePath);
        console.error(e);
        return null;
    }
}

function findDependencies(filePath, baseDir, testDir, visited = new Set()) {
    if (visited.has(filePath)) return;
    visited.add(filePath);
    console.log("visited",visited);

    const ast = parseFile(filePath);

    traverse(ast, {
        ImportDeclaration({ node }) {
            const importPath = node.source.value;
            const resolvedPath = resolveImportPath(importPath, filePath, baseDir, testDir);
            if (resolvedPath) {
                console.log("resolved path",resolvedPath);
                findDependencies(resolvedPath, baseDir,testDir, visited);
            }
        },
    });
}


function main(testFilePath, testDir) {
    const baseDir = path.dirname(testFilePath);
    const dependencies = new Set();
    findDependencies(testFilePath, baseDir,testDir, dependencies);
    dependencies.delete(testFilePath);
    console.log(JSON.stringify(Array.from(dependencies).map(dep => path.relative(projectRoot, dep))));
}
const args = process.argv.slice(2);
if (args.length !== 2) {
    console.error("Usage: node findTestDependencies.js <testFilePath> <projectRoot>");
    process.exit(1);
}

const [testFilePath, projectRoot] = args.map(arg => path.resolve(arg));
main(testFilePath, projectRoot);

// main(path.resolve('./src/foo.test.ts'), path.resolve('./src'));
